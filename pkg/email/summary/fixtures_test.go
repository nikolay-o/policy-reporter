package summary_test

import (
	"github.com/openreports/reports-api/pkg/client/clientset/versioned/fake"
	"github.com/openreports/reports-api/pkg/client/clientset/versioned/typed/openreports.io/v1alpha1"
	"go.uber.org/zap"

	"github.com/kyverno/policy-reporter/pkg/email"
	"github.com/kyverno/policy-reporter/pkg/validate"
)

var (
	filter = email.NewFilter(nil, validate.RuleSets{}, validate.RuleSets{})
	logger = zap.NewNop()
)

func NewFakeClient() (v1alpha1.OpenreportsV1alpha1Interface, v1alpha1.ReportInterface, v1alpha1.ClusterReportInterface) {
	client := fake.NewSimpleClientset().OpenreportsV1alpha1()

	return client, client.Reports("test"), client.ClusterReports()
}
