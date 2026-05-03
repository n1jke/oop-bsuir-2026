package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/domain"
)

func TestExchangeRequest_StateMachine(t *testing.T) {
	t.Run("Accept", func(t *testing.T) {
		tests := []struct {
			name           string
			initStatus     domain.ExchangeStatus
			wantErr        bool
			expectedStatus domain.ExchangeStatus
		}{
			{"valid", domain.Pending, false, domain.Accepted},
			{"invalid from accept", domain.Accepted, true, domain.Accepted},
			{"invalid from reject", domain.Rejected, true, domain.Rejected},
			{"invalid from complete", domain.Completed, true, domain.Completed},
			{"invalid from cancel", domain.Canceled, true, domain.Canceled},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				exchange := prepare(t, tt.initStatus)

				err := exchange.Accept()

				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}

				require.Equal(t, tt.expectedStatus, exchange.Status())
			})
		}
	})

	t.Run("Reject", func(t *testing.T) {
		tests := []struct {
			name           string
			initStatus     domain.ExchangeStatus
			wantErr        bool
			expectedStatus domain.ExchangeStatus
		}{
			{"valid", domain.Pending, false, domain.Rejected},
			{"invalid from accept", domain.Accepted, true, domain.Accepted},
			{"invalid from complete", domain.Completed, true, domain.Completed},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				exchange := prepare(t, tt.initStatus)

				err := exchange.Reject()

				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}

				require.Equal(t, tt.expectedStatus, exchange.Status())
			})
		}
	})

	t.Run("Complete", func(t *testing.T) {
		tests := []struct {
			name           string
			initStatus     domain.ExchangeStatus
			wantErr        bool
			expectedStatus domain.ExchangeStatus
		}{
			{"valid", domain.Accepted, false, domain.Completed},
			{"invalid from pending", domain.Pending, true, domain.Pending},
			{"invalid from reject", domain.Rejected, true, domain.Rejected},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				exchange := prepare(t, tt.initStatus)

				err := exchange.Complete()

				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}

				require.Equal(t, tt.expectedStatus, exchange.Status())
			})
		}
	})

	t.Run("Cancel", func(t *testing.T) {
		tests := []struct {
			name           string
			initStatus     domain.ExchangeStatus
			wantErr        bool
			expectedStatus domain.ExchangeStatus
		}{
			{"valid from pending", domain.Pending, false, domain.Canceled},
			{"valid from accepted", domain.Accepted, false, domain.Canceled},
			{"valid from rejected", domain.Rejected, false, domain.Canceled},
			{"invalid from complete", domain.Completed, true, domain.Completed},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				exchange := prepare(t, tt.initStatus)

				err := exchange.Cancel()

				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}

				require.Equal(t, tt.expectedStatus, exchange.Status())
			})
		}
	})
}

func newValidDatePolicy() *domain.DatePolicy {
	return &domain.DatePolicy{}
}

func prepare(t *testing.T, status domain.ExchangeStatus) *domain.ExchangeRequest {
	t.Helper()

	exchange, err := domain.NewExchangeRequest(uuid.New(), uuid.New(), uuid.New(), newValidDatePolicy(), "")
	require.NoError(t, err)

	switch status {
	case domain.Pending:
		return exchange
	case domain.Accepted:
		require.NoError(t, exchange.Accept())
		return exchange
	case domain.Rejected:
		require.NoError(t, exchange.Reject())
		return exchange
	case domain.Completed:
		require.NoError(t, exchange.Accept())
		require.NoError(t, exchange.Complete())

		return exchange
	case domain.Canceled:
		require.NoError(t, exchange.Cancel())
		return exchange
	default:
		t.Fatalf("unsupported status %v", status)
		return nil
	}
}
