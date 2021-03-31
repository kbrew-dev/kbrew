set -e

version=$(cut -d'=' -f2- .release)
if [[ -z ${version} ]]; then
    echo "Invalid version set in .release"
    exit 1
fi


if [[ -z ${GITHUB_TOKEN} ]]; then
    echo "GITHUB_TOKEN not set. Usage: GITHUB_TOKEN=<TOKEN> ./hack/release.sh"
    exit 1
fi

echo "Publishing release ${version}"

generate_changelog() {
    local version=$1

    # generate changelog from github
    github_changelog_generator --user kbrew-dev --project kbrew -t ${GITHUB_TOKEN} --future-release ${version} -o CHANGELOG.md
    sed -i '$d' CHANGELOG.md
}

git_tag() {
    local version=$1
    echo "Creating a git tag"
    git add .release CHANGELOG.md
    git commit -m "Release ${version}"
    git tag ${version}
    git push --tags origin main
    echo 'Git tag pushed successfully'
}

make_release() {
    goreleaser release --rm-dist --debug 
}

generate_changelog $version
git_tag $version
make_release

echo "=========================== Done ============================="
echo "Congratulations!! Release ${version} published."
echo "Don't forget to add changelog in the release description."
echo "=============================================================="
