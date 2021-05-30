package events

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/kbrew-dev/kbrew/pkg/config"
	"github.com/kbrew-dev/kbrew/pkg/kube"
	"github.com/kbrew-dev/kbrew/pkg/version"
)

// EventCatagory is the Google Analytics Event catagory
type EventCatagory string

const (
	kbrewTrackingID = "UA-195717361-1"
	gaCollectURL    = "https://www.google-analytics.com/collect"
	httpTimeout     = 5 * time.Second
)

var (
	k8sVersion string

	// ECInstallSuccess represents install success event catagory
	ECInstallSuccess EventCatagory = "install-success"
	// ECInstallFail represents install failure event catagory
	ECInstallFail EventCatagory = "install-fail"
	// ECInstallTimeout represents install timeout event catagory
	ECInstallTimeout EventCatagory = "install-timeout"
	// ECK8sEvent represents k8s events event catagory
	ECK8sEvent EventCatagory = "k8s-event"
)

type k8sEvent struct {
	Reason  string
	Message string
	Object  string
	Action  string
}

// KbrewEvent contains information to report Event to Google Analytics
type KbrewEvent struct {
	gaVersion    string
	gaType       string
	gaTID        string
	gaCID        string
	gaAIP        string
	gaAppName    string
	gaAppVersion string
	gaEvCatagory string
	gaEvAction   string
	gaEvLabel    string
	gaKbrewArgs  string
}

// String returns string representation of Event Catagory
func (ec EventCatagory) String() string {
	return string(ec)
}

// NewKbrewEvent return new KbrewEvent
func NewKbrewEvent(appConfig *config.AppConfig) *KbrewEvent {
	return &KbrewEvent{
		gaVersion:    "1",
		gaType:       "event",
		gaTID:        kbrewTrackingID,
		gaCID:        viper.GetString(config.AnalyticsUUID),
		gaAIP:        "1",
		gaAppName:    "kbrew",
		gaAppVersion: version.Short(),
		gaEvLabel:    fmt.Sprintf("k8s %s", k8sVersion),
		gaEvAction:   appConfig.App.Name,
		gaKbrewArgs:  labels.FormatLabels(argsToLabels(appConfig.App.Args)),
	}
}

// Report sends event to Google Analytics
func (kv *KbrewEvent) Report(ctx context.Context, ec EventCatagory, err error, k8sEvent *k8sEvent) error {
	v := url.Values{
		"v":   {kv.gaVersion},
		"tid": {kv.gaTID},
		"cid": {kv.gaCID},
		"aip": {kv.gaAIP},
		"t":   {kv.gaType},
		"ec":  {ec.String()},
		"ea":  {kv.gaEvAction},
		"el":  {kv.gaEvLabel},
		"an":  {kv.gaAppName},
		"av":  {kv.gaAppVersion},
		"cd1": {},
		"cd2": {},
		"cd3": {},
		"cd4": {},
		"cd5": {},
		"cd6": {kv.gaKbrewArgs},
	}

	if err != nil {
		// Set kbrew message
		v.Set("cd5", err.Error())
	}

	if k8sEvent != nil {
		// Set k8s_reason
		v.Set("cd1", k8sEvent.Reason)
		// Set k8s_message
		v.Set("cd2", k8sEvent.Message)
		// Set k8s_action
		v.Set("cd3", k8sEvent.Action)
		// Set k8s_object
		v.Set("cd4", k8sEvent.Object)
	}

	buf := bytes.NewBufferString(v.Encode())
	req, err1 := http.NewRequest("POST", gaCollectURL, buf)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", fmt.Sprintf("kbrew/%s", version.Short()))
	if err1 != nil {
		return err1
	}
	ctx, cancel := context.WithTimeout(ctx, httpTimeout)
	defer cancel()

	req = req.WithContext(ctx)

	client := http.DefaultClient
	resp, err1 := client.Do(req)
	if err1 != nil {
		return err1
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Analytics report failed with status code %d", resp.StatusCode))
	}
	defer resp.Body.Close()
	return err1
}

// ReportK8sEvents sends kbrew events with K8s events to Google Analytics
func (kv *KbrewEvent) ReportK8sEvents(ctx context.Context, err error, workloads []corev1.ObjectReference) error {
	k8sEvents, err1 := getPodEvents(ctx, workloads)
	if err1 != nil {
		return err1
	}
	for _, event := range k8sEvents {
		err1 := kv.Report(ctx, ECK8sEvent, err, &event)
		if err1 != nil {
			return err1
		}
	}
	return nil
}

func getPodEvents(ctx context.Context, workloads []corev1.ObjectReference) ([]k8sEvent, error) {
	notRunningPods, err := kube.FetchNonRunningPods(ctx, workloads)
	if err != nil {
		return nil, err
	}
	events := []k8sEvent{}
	for _, pod := range notRunningPods {
		ks8Events, err := getK8sEvents(ctx, corev1.ObjectReference{Name: pod.GetName(), Namespace: pod.GetNamespace(), UID: pod.GetUID(), Kind: "Pod"})
		if err != nil {
			return nil, err
		}
		events = append(events, ks8Events...)
	}
	return events, nil
}

func prepareObjectSelector(objReference corev1.ObjectReference) string {
	return labels.Set{
		"involvedObject.name":      objReference.Name,
		"involvedObject.namespace": objReference.Namespace,
		"involvedObject.uid":       string(objReference.UID),
		"involvedObject.kind":      objReference.Kind,
		"type":                     "Warning",
	}.String()
}

func init() {
	var err error
	k8sVersion, err = getK8sVersion()
	if err != nil {
		fmt.Printf("ERROR: Failed to get K8s version. %s", err.Error)
	}
}

func getK8sEvents(ctx context.Context, objReference corev1.ObjectReference) ([]k8sEvent, error) {
	clis, err := kube.NewClient()
	if err != nil {
		return nil, err
	}
	objSelector := prepareObjectSelector(objReference)
	eventList, err := clis.KubeCli.CoreV1().Events(objReference.Namespace).List(ctx, metav1.ListOptions{FieldSelector: objSelector})
	if err != nil {
		return nil, err
	}
	retEventList := []k8sEvent{}
	for _, event := range eventList.Items {
		objRef := corev1.ObjectReference{
			Name:      event.InvolvedObject.Name,
			Namespace: event.InvolvedObject.Namespace,
			Kind:      event.InvolvedObject.Kind,
		}
		retEventList = append(retEventList, k8sEvent{
			Reason:  event.Reason,
			Message: event.Message,
			Object:  objRef.String(),
			Action:  event.Action,
		})
	}
	return retEventList, nil
}

func getK8sVersion() (string, error) {
	clis, err := kube.NewClient()
	if err != nil {
		return "", err
	}
	versionInfo, err := clis.DiscoveryCli.ServerVersion()
	if err != nil {
		return "", err
	}
	return versionInfo.String(), nil
}

func argsToLabels(args map[string]interface{}) map[string]string {
	labels := make(map[string]string)
	for k, v := range args {
		labels[k] = fmt.Sprintf("%v", v)
	}
	return labels
}
