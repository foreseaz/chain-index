package httpapi_test

import (
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/crypto-com/chainindex/adapter"
	"github.com/crypto-com/chainindex/adapter/httpapi"
	. "github.com/crypto-com/chainindex/adapter/httpapi/test"
	. "github.com/crypto-com/chainindex/adapter/httpapi/test/mock"
	. "github.com/crypto-com/chainindex/usecase/test/fake"
	"github.com/crypto-com/chainindex/usecase/viewrepo"
	. "github.com/crypto-com/chainindex/usecase/viewrepo/test/mock"
)

var _ = Describe("CouncilNodes", func() {
	var mockCouncilNodeViewRepo *MockCouncilNodeViewRepo
	var mockRoutePath *MockRoutePath
	var mockHandler *httpapi.CouncilNodesHandler

	BeforeEach(func() {
		fakeLogger := &FakeLogger{}
		mockCouncilNodeViewRepo = &MockCouncilNodeViewRepo{}
		mockRoutePath = &MockRoutePath{}

		mockHandler = httpapi.NewCouncilNodesHandler(fakeLogger, mockRoutePath, mockCouncilNodeViewRepo)
	})

	Describe("ListActiveCouncilNodes", func() {
		It("should return BadRequest when pagination is missing", func() {
			reqWithInvalidPage := NewMockHTTPGetRequest(HTTPQueryParams{
				"page": "invalid",
			})
			respSpy := httptest.NewRecorder()

			mockHandler.ListActiveCouncilNodes(respSpy, reqWithInvalidPage)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})
	})

	Describe("FindCouncilNodeById", func() {
		It("should return BadRequest when id is missing", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{})

			mockHandler.FindCouncilNodeById(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when id has invalid type", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"id": "invalid",
			})

			mockHandler.FindCouncilNodeById(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when id is negative number", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"id": "-10",
			})

			mockHandler.FindCouncilNodeById(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when id is decimal number", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"id": "1.1",
			})

			mockHandler.FindCouncilNodeById(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return NotFound when council node does not exist", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"id": "10",
			})

			mockCouncilNodeViewRepo.On(
				"FindById", mock.Anything,
			).Return((*viewrepo.CouncilNode)(nil), adapter.ErrNotFound)

			mockHandler.FindCouncilNodeById(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(404))
		})
	})

	Describe("ListCouncilNodeActivitiesById", func() {
		It("should return BadRequest when id is missing", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{})

			mockHandler.ListCouncilNodeActivitiesById(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when id has invalid type", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"id": "invalid",
			})

			mockHandler.ListCouncilNodeActivitiesById(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when id is negative number", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"id": "-10",
			})

			mockHandler.ListCouncilNodeActivitiesById(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when id is decimal number", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"id": "1.1",
			})

			mockHandler.ListCouncilNodeActivitiesById(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return NotFound when council node does not exist", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"id": "10",
			})

			mockCouncilNodeViewRepo.On(
				"ListActivitiesById", mock.Anything, mock.Anything,
			).Return(([]viewrepo.StakingAccountActivity)(nil), (*viewrepo.PaginationResult)(nil), adapter.ErrNotFound)

			mockHandler.ListCouncilNodeActivitiesById(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(404))
		})
	})
})