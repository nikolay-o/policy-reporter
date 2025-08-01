package metrics

import (
	"github.com/kyverno/policy-reporter/pkg/openreports"
	"github.com/kyverno/policy-reporter/pkg/report"
	"github.com/kyverno/policy-reporter/pkg/validate"
)

func NewResultFilter(namespace, status, policy, source, severity, kind validate.RuleSets) *report.ResultFilter {
	f := &report.ResultFilter{}
	if namespace.Count() > 0 {
		f.AddValidation(func(r openreports.ResultAdapter) bool {
			if !r.HasResource() {
				return true
			}

			return validate.Namespace(r.GetResource().Namespace, namespace)
		})
	}

	if status.Count() > 0 {
		f.AddValidation(func(r openreports.ResultAdapter) bool {
			return validate.MatchRuleSet(string(r.Result), status)
		})
	}

	if policy.Count() > 0 {
		f.AddValidation(func(r openreports.ResultAdapter) bool {
			return validate.MatchRuleSet(r.Policy, policy)
		})
	}

	if source.Count() > 0 {
		f.AddValidation(func(r openreports.ResultAdapter) bool {
			return validate.MatchRuleSet(r.Source, source)
		})
	}

	if severity.Count() > 0 {
		f.AddValidation(func(r openreports.ResultAdapter) bool {
			return validate.MatchRuleSet(string(r.Severity), severity)
		})
	}

	if kind.Count() > 0 {
		f.AddValidation(func(r openreports.ResultAdapter) bool {
			if !r.HasResource() {
				return true
			}

			return validate.Kind(r.GetResource().Kind, kind)
		})
	}

	return f
}

func NewReportFilter(namespace, source validate.RuleSets) *report.ReportFilter {
	f := &report.ReportFilter{}
	if namespace.Count() > 0 {
		f.AddValidation(func(r openreports.ReportInterface) bool {
			return validate.Namespace(r.GetNamespace(), namespace)
		})
	}

	if source.Count() > 0 {
		f.AddValidation(func(r openreports.ReportInterface) bool {
			if r.GetSource() == "" {
				return true
			}

			return validate.MatchRuleSet(r.GetSource(), source)
		})
	}

	return f
}
