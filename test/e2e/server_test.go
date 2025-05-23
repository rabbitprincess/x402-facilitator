package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rabbitprincess/x402-facilitator/api"
	"github.com/rabbitprincess/x402-facilitator/test/mock"
	"github.com/rabbitprincess/x402-facilitator/types"
)

var _ = Describe("Server", func() {
	var (
		server http.Handler
		resp   *http.Response
	)

	BeforeEach(func() {
		// Create a new server instance with mock facilitator
		server = api.NewServer(&mock.Facilitator{})
	})

	Context("when accessing an undefined route", func() {
		BeforeEach(func() {
			// Create a new HTTP request
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			// Create a response recorder
			w := httptest.NewRecorder()

			// Serve the request
			server.ServeHTTP(w, req)

			// Get the response
			resp = w.Result()
		})

		It("should return 404 Not Found", func() {
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})
	})

	Context("when accessing /verify endpoint", func() {
		var (
			payload    []byte
			statusCode int
			respBody   map[string]interface{}
		)

		JustBeforeEach(func() {
			// Create a new HTTP request
			req := httptest.NewRequest(http.MethodPost, "/verify", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			w := httptest.NewRecorder()

			// Serve the request
			server.ServeHTTP(w, req)

			// Get the response
			resp = w.Result()
			statusCode = resp.StatusCode

			// Parse response body
			if statusCode != http.StatusNoContent {
				err := json.NewDecoder(resp.Body).Decode(&respBody)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		Context("with valid request", func() {
			BeforeEach(func() {
				req := types.PaymentVerifyRequest{
					X402Version:   1,
					PaymentHeader: `{"x402Version":1,"scheme":"evm","network":"base-sepolia","payload":{}}`,
					PaymentRequirements: types.PaymentRequirements{
						Scheme:            "evm",
						Network:           "base-sepolia",
						MaxAmountRequired: "1000000000000000000",
						Resource:          "/api/resource",
						Description:       "Test resource",
						MimeType:          "application/json",
						PayTo:             "0x1234567890123456789012345678901234567890",
						MaxTimeoutSeconds: 30,
						Asset:             "0x1234567890123456789012345678901234567890",
					},
				}
				var err error
				payload, err = json.Marshal(req)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return 400 Bad Request", func() {
				Expect(statusCode).To(Equal(http.StatusBadRequest))
				Expect(respBody).To(HaveKey("message"))
			})
		})

		Context("with invalid request (missing required fields)", func() {
			BeforeEach(func() {
				// Deliberately omit required fields
				req := map[string]interface{}{
					"x402Version": 1,
					// Missing paymentHeader and paymentRequirements
				}
				var err error
				payload, err = json.Marshal(req)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return 400 Bad Request", func() {
				Expect(statusCode).To(Equal(http.StatusBadRequest))
				Expect(respBody).To(HaveKey("message"))
			})
		})
	})

	Context("when accessing /settle endpoint", func() {
		var (
			payload    []byte
			statusCode int
			respBody   map[string]interface{}
		)

		JustBeforeEach(func() {
			// Create a new HTTP request
			req := httptest.NewRequest(http.MethodPost, "/settle", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			w := httptest.NewRecorder()

			// Serve the request
			server.ServeHTTP(w, req)

			// Get the response
			resp = w.Result()
			statusCode = resp.StatusCode

			// Parse response body
			if statusCode != http.StatusNoContent {
				err := json.NewDecoder(resp.Body).Decode(&respBody)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		Context("with valid request", func() {
			BeforeEach(func() {
				req := types.PaymentSettleRequest{
					X402Version:   1,
					PaymentHeader: `{"x402Version":1,"scheme":"evm","network":"base-sepolia","payload":{}}`,
					PaymentRequirements: types.PaymentRequirements{
						Scheme:            "evm",
						Network:           "base-sepolia",
						MaxAmountRequired: "1000000000000000000",
						Resource:          "/api/resource",
						Description:       "Test resource",
						MimeType:          "application/json",
						PayTo:             "0x1234567890123456789012345678901234567890",
						MaxTimeoutSeconds: 30,
						Asset:             "0x1234567890123456789012345678901234567890",
					},
				}
				var err error
				payload, err = json.Marshal(req)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return success response", func() {
				Expect(statusCode).To(Equal(http.StatusNotImplemented))
				Expect(respBody).To(HaveKey("message"))
				Expect(respBody["message"]).To(Equal("settlement not implemented yet"))
			})
		})

		Context("with invalid request (missing required fields)", func() {
			BeforeEach(func() {
				// Deliberately omit required fields
				req := map[string]interface{}{
					"x402Version": 1,
					// Missing paymentHeader and paymentRequirements
				}
				var err error
				payload, err = json.Marshal(req)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return 501 Not Implemented", func() {
				Expect(statusCode).To(Equal(http.StatusNotImplemented))
				Expect(respBody).To(HaveKey("message"))
				Expect(respBody["message"]).To(Equal("settlement not implemented yet"))
			})
		})
	})
})
