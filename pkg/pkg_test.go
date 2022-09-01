package pkg

// func TestToolchainGetCatalog(t *testing.T) {
// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte(`{
//   "hookSource":"quay.io/trustacks/test:latest",
//   "components":{
//     "test":{
//       "repository":"https://test-charts.trustacks.io",
//       "chart":"test/test",
//       "version":"1.1.1"
//     }
//   }
// }`))
// 	}))
// 	catalog, err := GetCatalog(ts.URL)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if catalog.HookSource != "quay.io/trustacks/test:latest" {
// 		t.Fatal("got an unexpected hook source")
// 	}
// 	if catalog.Components["test"].Repo != "https://test-charts.trustacks.io" {
// 		t.Fatal("got an unexpected repo")
// 	}
// 	if catalog.Components["test"].Chart != "test/test" {
// 		t.Fatal("got an unexpected chart")
// 	}
// 	if catalog.Components["test"].Version != "1.1.1" {
// 		t.Fatal("got an unexpected chart")
// 	}
// }

// func TestToolchainAddComponents(t *testing.T) {
// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		switch r.URL.Path {
// 		case "/charts/helloworld-1.0.0.tgz":
// 			data, err := ioutil.ReadFile("testdata/helloworld-1.0.0.tgz")
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			_, err = w.Write(data)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 		}
// 	}))
// 	catalog := &ComponentCatalog{
// 		Components: map[string]Component{
// 			"helloworld": {
// 				Repo:    fmt.Sprintf("%s/charts", ts.URL),
// 				Chart:   "helloworld",
// 				Version: "1.0.0",
// 				Hooks:   "",
// 			},
// 		},
// 	}
// 	d, err := ioutil.TempDir("", "test")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if err := AddComponents(d, []string{"helloworld"}, catalog); err != nil {
// 		t.Fatal(err)
// 	}
// 	_, err = os.Stat(fmt.Sprintf("%s/components/helloworld", d))
// 	if os.IsNotExist(err) {
// 		t.Fatal("expected chart to exist")
// 	}
// }

// func TestAddHooks(t *testing.T) {
// 	hooksManifest, err := ioutil.ReadFile("testdata/hooks.yaml")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	catalog := &ComponentCatalog{
// 		HookSource: "quay.io/trustacks/test-catalog:latest",
// 		Components: map[string]Component{
// 			"helloworld": {
// 				Hooks: string(hooksManifest),
// 			},
// 		},
// 	}
// 	d, err := ioutil.TempDir("", "test")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if err := os.MkdirAll(fmt.Sprintf("%s/components", d), 0755); err != nil {
// 		t.Fatal(err)
// 	}
// 	cmd := exec.Command("cp", "-R", "testdata/helloworld", fmt.Sprintf("%s/components/helloworld", d))
// 	cmd.Stderr = os.Stderr
// 	if err := cmd.Run(); err != nil {
// 		t.Fatal(err)
// 	}
// 	if err := AddHooks(d, []string{"helloworld"}, catalog, map[string]interface{}{"testParam": "test"}); err != nil {
// 		t.Fatal(err)
// 	}
// 	_, err = os.Stat(fmt.Sprintf("%s/components/helloworld/templates/trustacks-hooks.yaml", d))
// 	if os.IsNotExist(err) {
// 		t.Fatal("expected hooks manifest to exist")
// 	}
// }

// func TestToolchainAddSubchartValues(t *testing.T) {
// 	catalog := &ComponentCatalog{
// 		HookSource: "quay.io/trustacks/test-catalog:latest",
// 		Components: map[string]Component{
// 			"helloworld": {
// 				Values: `username: username
// password: password`,
// 			},
// 		},
// 	}
// 	parameters := map[string]interface{}{
// 		"username": "username",
// 		"password": "password",
// 	}
// 	d, err := ioutil.TempDir("", "test")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if err := os.MkdirAll(path.Join(d, "components", "helloworld"), 0755); err != nil {
// 		t.Fatal(err)
// 	}
// 	if err := AddSubChartValues(d, []string{"helloworld"}, catalog, parameters); err != nil {
// 		t.Fatal(err)
// 	}
// 	values, err := ioutil.ReadFile(path.Join(d, "components", "helloworld", "override-values.yaml"))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	expectedValues := `username: username
// password: password`

// 	if string(values) != expectedValues {
// 		t.Fatal("got an unexpected values output")
// 	}
// }

// func TestRender(t *testing.T) {
// 	raw := `concourse:
// 	web:
// 		localAuth:
// 		enabled: true
// 		auth:
// 		mainTeam:
// 			oidc:
// 			group: admins
// 		oidc:
// 			enabled: true
// 			displayName: sso
// 			{{- if eq .sso "authentik"}}
// 			issuer: "{{- if eq .tls true -}}https{{- else -}}http{{- end -}}://authentik.{{ .domain }}:{{ .ingressPort }}/application/o/concourse"
// 			{{- end }}
// 		externalUrl: "{{- if eq .tls true -}}https{{- else -}}http{{- end -}}://concourse.{{ .domain }}:{{ .ingressPort }}"
// 		kubernetes:
// 		namespacePrefix: workflow-
// 		keepNamespace: false
// 	web:
// 	ingress:
// 		enabled: true
// 		hosts:
// 		- concourse.{{ .domain }}
// 	worker:
// 	env:
// 	- name: CONCOURSE_GARDEN_ALLOW_HOST_ACCESS
// 		value: "true"
// 	fullnameOverride: concourse
// 	postgresql:
// 	fullnameOverride: concourse-postgresql`
// 	tpl := template.Must(template.New("test").Funcs(sprig.FuncMap()).Parse(raw))
// 	var buf bytes.Buffer
// 	if err := tpl.Execute(&buf, map[string]interface{}{"domain": "local.gd", "ingressPort": "8081", "sso": "authentik", "tls": true}); err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Println(buf.String())
// }
