package domain_test

import (
	"testing"

	"github.com/davidgrcias/digital-wallet/internal/domain"
)

func TestWithdrawRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     domain.WithdrawRequest
		wantErr bool
	}{
		{
			name: "Valid Amount",
			req: domain.WithdrawRequest{
				Amount:      50000,
				Description: "Jajan",
			},
			wantErr: false,
		},
		{
			name: "Invalid Amount - Zero",
			req: domain.WithdrawRequest{
				Amount: 0,
			},
			wantErr: true,
		},
		{
			name: "Invalid Amount - Negative",
			req: domain.WithdrawRequest{
				Amount: -1000,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
