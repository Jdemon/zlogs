package main

import (
	"context"

	"github.com/Jdemon/zlogs"
)

func main() {

	zlogs.NewLogger(&zlogs.Config{
		Level: "debug",
		Masking: zlogs.ConfigMasking{
			Enabled:         true,
			SensitiveFields: []string{"lastName"},
		},
	})

	data := map[string]any{
		"data": map[string]any{
			"password":      "P@ssw0rd",
			"mobile_number": "0909263742",
			"id":            "112132321312",
			"firstname":     "John",
			"lastName":      "Doe",
		},
		"credit_card": "4231234512341234",
	}

	ctx := context.WithValue(context.Background(), zlogs.TraceID, "trace-id-value")
	event := zlogs.Info().WithFields(data).Ctx(ctx)
	event2 := zlogs.Info().WithFields(data)
	event.Msgf("hello world %s", "jay")
	event2.Msgf("hello world %s", "jay2")

	zlogs.NewGORMLogger(&zlogs.Config{
		Level: "debug",
		Masking: zlogs.ConfigMasking{
			Enabled:         true,
			SensitiveFields: []string{"lastName"},
		},
	}).Error(ctx, "hello gorm orm")

	zlogs.Debug().Msg("test")
}
