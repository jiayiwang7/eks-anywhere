package v1alpha1

import (
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetSnowDatacenterConfig(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		want     *SnowDatacenterConfig
		wantErr  string
	}{
		{
			name:     "file doesn't exist",
			fileName: "testdata/fake_file.yaml",
			want:     nil,
			wantErr:  "no such file or directory",
		},
		{
			name:     "not parseable file",
			fileName: "testdata/not_parseable_cluster_snow.yaml",
			want:     nil,
			wantErr:  "error unmarshaling JSON: while decoding JSON: json: unknown field",
		},
		{
			name:     "valid 1.21",
			fileName: "testdata/cluster_1_21_snow.yaml",
			want: &SnowDatacenterConfig{
				TypeMeta: metav1.TypeMeta{
					Kind:       SnowDatacenterKind,
					APIVersion: SchemeBuilder.GroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "eksa-unit-test",
				},
				Spec: SnowDatacenterConfigSpec{},
			},
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			got, err := GetSnowDatacenterConfig(tt.fileName)
			if tt.wantErr == "" {
				g.Expect(err).To(BeNil())
			} else {
				g.Expect(err).To(MatchError(ContainSubstring(tt.wantErr)))
			}

			g.Expect(got).To(Equal(tt.want))
		})
	}
}

func TestSnowDatacenterConfigSetDefaults(t *testing.T) {
	tests := []struct {
		name   string
		before *SnowDatacenterConfig
		after  *SnowDatacenterConfig
	}{
		{
			name: "identity ref nil",
			before: &SnowDatacenterConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: SnowDatacenterConfigSpec{},
			},
			after: &SnowDatacenterConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: SnowDatacenterConfigSpec{
					IdentityRef: &Ref{
						Name: "test-snow-credentials",
						Kind: "Secret",
					},
				},
			},
		},
		{
			name: "identity ref exists",
			before: &SnowDatacenterConfig{
				Spec: SnowDatacenterConfigSpec{
					IdentityRef: &Ref{
						Name: "creds-1",
						Kind: "Secret",
					},
				},
			},
			after: &SnowDatacenterConfig{
				Spec: SnowDatacenterConfigSpec{
					IdentityRef: &Ref{
						Name: "creds-1",
						Kind: "Secret",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			tt.before.SetDefaults()
			g.Expect(tt.before).To(Equal(tt.after))
		})
	}
}

func TestSnowDatacenterConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		obj     *SnowDatacenterConfig
		wantErr string
	}{
		{
			name: "identity nil",
			obj: &SnowDatacenterConfig{
				Spec: SnowDatacenterConfigSpec{},
			},
			wantErr: "",
		},
		{
			name: "valid identity ref",
			obj: &SnowDatacenterConfig{
				Spec: SnowDatacenterConfigSpec{
					IdentityRef: &Ref{
						Name: "creds-1",
						Kind: "Secret",
					},
				},
			},
			wantErr: "",
		},
		{
			name: "invalid identity ref kind",
			obj: &SnowDatacenterConfig{
				Spec: SnowDatacenterConfigSpec{
					IdentityRef: &Ref{
						Name: "creds-1",
						Kind: "UnknowKind",
					},
				},
			},
			wantErr: "SnowDatacenterConfig IdetityRef Kind UnknowKind is invalid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			err := tt.obj.Validate()
			if tt.wantErr == "" {
				g.Expect(err).To(BeNil())
			} else {
				g.Expect(err).To(MatchError(ContainSubstring(tt.wantErr)))
			}
		})
	}
}
