/*
Copyright 2017 The Kedge Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package spec

import (
	"reflect"
	"testing"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	api_v1 "k8s.io/kubernetes/pkg/api/v1"
)

func TestUnmarshalValidateFixApp(t *testing.T) {
	tests := []struct {
		Name string
		Data []byte
		App  *App
	}{
		{
			Name: "One container mentioned in the spec",
			Data: []byte(`
name: test
deployments:
  - containers:
    - image: nginx
services:
  - ports:
    - port: 8080`),
			App: &App{

				ObjectMeta: meta_v1.ObjectMeta{
					Name: "test",
					Labels: map[string]string{
						appLabelKey: "test",
					},
				},
				Deployments: []DeploymentSpecMod{
					{
						ObjectMeta: meta_v1.ObjectMeta{
							Name: "test",
							Labels: map[string]string{
								appLabelKey: "test",
							},
						},
						PodSpecMod: PodSpecMod{
							Containers: []Container{{Container: api_v1.Container{Name: "test", Image: "nginx"}}},
						},
					},
				},
				Services: []ServiceSpecMod{
					{
						ObjectMeta: meta_v1.ObjectMeta{
							Name: "test",
							Labels: map[string]string{
								appLabelKey: "test",
							},
						},
						Ports: []ServicePortMod{{ServicePort: api_v1.ServicePort{Port: 8080}}}},
				},
			},
		},
		{
			Name: "One persistent volume mentioned in the spec",
			Data: []byte(`
name: test
deployments:
  - containers:
    - image: nginx
services:
  - ports:
    - port: 8080
volumeClaims:
- size: 500Mi`),
			App: &App{
				ObjectMeta: meta_v1.ObjectMeta{
					Name: "test",
					Labels: map[string]string{
						appLabelKey: "test",
					},
				},
				Deployments: []DeploymentSpecMod{
					{
						ObjectMeta: meta_v1.ObjectMeta{
							Name: "test",
							Labels: map[string]string{
								appLabelKey: "test",
							},
						},
						PodSpecMod: PodSpecMod{
							Containers: []Container{{Container: api_v1.Container{Name: "test", Image: "nginx"}}},
						},
					},
				},
				Services: []ServiceSpecMod{
					{
						ObjectMeta: meta_v1.ObjectMeta{
							Name: "test",
							Labels: map[string]string{
								appLabelKey: "test",
							},
						},
						Ports: []ServicePortMod{
							{
								ServicePort: api_v1.ServicePort{
									Port: 8080,
								},
							},
						},
					},
				},
				VolumeClaims: []VolumeClaim{
					{
						ObjectMeta: meta_v1.ObjectMeta{
							Name: "test",
							Labels: map[string]string{
								appLabelKey: "test",
							},
						},
						Size: "500Mi"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var app App

			if err := app.LoadData(test.Data); err != nil {
				t.Fatalf("unable to unmarshal data - %v", err)
			}

			if err := app.Validate(); err != nil {
				t.Fatalf("unable to validate data - %v", err)
			}

			if err := app.Fix(); err != nil {
				t.Fatalf("unable to fix data - %v", err)
			}

			if !reflect.DeepEqual(test.App, &app) {
				t.Fatalf("==> Expected:\n%v\n==> Got:\n%v", PrettyPrintObjects(test.App), PrettyPrintObjects(app))
			}
		})
	}
}
