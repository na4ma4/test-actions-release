package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
	"time"

	ts "github.com/na4ma4/go-timestring"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// basicFunctions are the set of initial functions provided to every template.
func basicFunctions(extra ...template.FuncMap) template.FuncMap {
	basicFuncMap := template.FuncMap{
		"json": func(v interface{}) string {
			buf := &bytes.Buffer{}
			enc := json.NewEncoder(buf)
			enc.SetEscapeHTML(false)
			_ = enc.Encode(v) //nolint:errchkjson // templating function
			// Remove the trailing new line added by the encoder
			return strings.TrimSpace(buf.String())
		},
		"split":    strings.Split,
		"join":     strings.Join,
		"title":    cases.Title(language.English).String,
		"lower":    cases.Lower(language.English).String,
		"upper":    cases.Upper(language.English).String,
		"pad":      padWithSpace,
		"padlen":   padToLength,
		"padmax":   padToMaxLength,
		"truncate": truncateWithLength,
		"tf":       stringTrueFalse,
		"yn":       stringYesNo,
		"t":        stringTab,
		"age":      humanAgeFormat,
		"time":     timeFormat,
		"date":     dateFormat,
	}

	if len(extra) > 0 {
		for _, add := range extra {
			for k, v := range add {
				basicFuncMap[k] = v
			}
		}
	}

	return basicFuncMap
}

// padToLength adds whitespace to pad to the supplied length.
func padToMaxLength(source interface{}) string {
	switch val := source.(type) {
	case *string:
		return padToLength(*val, 0)
	case *int:
		return padToLength(fmt.Sprintf("%d", *val), 0)
	case *int8:
		return padToLength(fmt.Sprintf("%d", *val), 0)
	case *int16:
		return padToLength(fmt.Sprintf("%d", *val), 0)
	case *int32:
		return padToLength(fmt.Sprintf("%d", *val), 0)
	case *int64:
		return padToLength(fmt.Sprintf("%d", *val), 0)
	default:
		return padToLength(source, 0)
	}
}

// padToLength adds whitespace to pad to the supplied length.
func padToLength(source interface{}, prefix int) string {
	switch val := source.(type) {
	case *string:
		return padToLength(*val, prefix)
	case *int:
		return padToLength(fmt.Sprintf("%d", *val), prefix)
	case *int8:
		return padToLength(fmt.Sprintf("%d", *val), prefix)
	case *int16:
		return padToLength(fmt.Sprintf("%d", *val), prefix)
	case *int32:
		return padToLength(fmt.Sprintf("%d", *val), prefix)
	case *int64:
		return padToLength(fmt.Sprintf("%d", *val), prefix)
	default:
		return fmt.Sprintf(fmt.Sprintf("%%-%ds", prefix), source)
	}
}

// padWithSpace adds whitespace to the input if the input is non-empty.
func padWithSpace(source interface{}, prefix, suffix int) string {
	src := fmt.Sprintf("%s", source)

	if src == "" {
		return src
	}

	return strings.Repeat(" ", prefix) + src + strings.Repeat(" ", suffix)
}

// humanAgeFormat returns a duration in a human readable format.
func humanAgeFormat(source interface{}) string {
	switch src := source.(type) {
	case time.Time:
		return ts.LongProcess.Option(ts.Abbreviated, ts.ShowMSOnSeconds).String(time.Since(src))
	case timestamppb.Timestamp:
		return ts.LongProcess.Option(ts.Abbreviated, ts.ShowMSOnSeconds).String(time.Since(src.AsTime()))
	case *timestamppb.Timestamp:
		return ts.LongProcess.Option(ts.Abbreviated, ts.ShowMSOnSeconds).String(time.Since(src.AsTime()))
	default:
		return fmt.Sprintf("%s", src)
	}
}

// // ageFormat returns time in seconds ago.
// func ageFormat(source interface{}) string {
// 	switch s := source.(type) {
// 	case time.Time:
// 		return fmt.Sprintf("%0.2f", time.Since(s).Seconds())
// 		// return s.Format(time.RFC3339)
// 	case timestamppb.Timestamp:
// 		return fmt.Sprintf("%0.2f", time.Since(s.AsTime()).Seconds())
// 	case *timestamppb.Timestamp:
// 		return fmt.Sprintf("%0.2f", time.Since(s.AsTime()).Seconds())
// 	default:
// 		return fmt.Sprintf("%s", source)
// 	}
// }

// timeFormat returns time in RFC3339 format.
func timeFormat(source interface{}) string {
	switch src := source.(type) {
	case time.Time:
		return src.Format(time.RFC3339)
	case timestamppb.Timestamp:
		return src.AsTime().Format(time.RFC3339)
	case *timestamppb.Timestamp:
		return src.AsTime().Format(time.RFC3339)
	default:
		return fmt.Sprintf("%s", src)
	}
}

// dateFormat returns date in YYYY-MM-DD format.
func dateFormat(source interface{}) string {
	switch src := source.(type) {
	case time.Time:
		return src.Format("2006-01-02")
	case timestamppb.Timestamp:
		return src.AsTime().Format("2006-01-02")
	case *timestamppb.Timestamp:
		return src.AsTime().Format("2006-01-02")
	default:
		return fmt.Sprintf("%q", src)
	}
}

func stringBool(source interface{}, trueString, falseString string) string {
	switch val := source.(type) {
	case *bool:
		if val != nil {
			return stringBool(*val, trueString, falseString)
		}

		return "nil"
	case bool:
		if val {
			return trueString
		}

		return falseString
	default:
		return fmt.Sprintf("%s", val)
	}
}

// stringTrueFalse returns "true" or "false" for boolean input.
func stringTrueFalse(source interface{}) string {
	return stringBool(source, "true", "false")
}

// stringYesNo returns "yes" or "no" for boolean input.
func stringYesNo(source bool) string {
	return stringBool(source, "yes", "no")
}

// stringTab returns a tab character.
func stringTab() string {
	return "\t"
}

// truncateWithLength truncates the source string up to the length provided by the input.
func truncateWithLength(source interface{}, length int) string {
	src := fmt.Sprintf("%s", source)

	if len(src) < length {
		return src
	}

	return src[:length]
}
