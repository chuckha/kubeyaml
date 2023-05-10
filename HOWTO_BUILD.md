# Instructions for building kubeyaml

```sh
k8s_version=1.25
curl -L https://github.com/kubernetes/kubernetes/blob/release-$k8s_version/api/openapi-spec/swagger.json -o backend/internal/kubernetes/data/swagger-$k8s_version.json
// update releases in `backend/scripts/update-schemas.go` and in `backend/internal/kubernetes/data/lookup.go`
// update versions in `backend/cmd/kubeyaml/kubeyaml.go`
rm -f backend/internal/service/validation/data/swagger-*.json
cp backend/internal/kubernetes/data/swagger-*.json backend/internal/service/validation/data/
go build ./backend/scripts/update-schemas.go
mv update-schemas ./backend/scripts/
cd backend/
./scripts/update-schemas

// test
cd -
cd backend/
go build -a -installsuffix cgo -o kubeyaml ./cmd/server
./kubeyaml --help
go build -a -installsuffix cgo -o kubeyaml ./cmd/kubeyaml
./kubeyaml --help

cd -
rm -f backend/scripts/update-schemas backend/kubeyaml
git add backend/internal/kubernetes backend/internal/service/validation
git commit -a -m "add versions"
git push

// get new commit 
go get github.com/cristifalcas/kubeyaml/backend@1f0520f8d81a0fbcd26542cef87716eb2380f2e3
./dev-scripts/dependencies/generate.sh --go
bazel run @com_github_cristifalcas_kubeyaml_backend//cmd/kubeyaml -- --help
```
