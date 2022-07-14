package main_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"taskism/controllers"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Backend", Ordered, func() {
	var router *gin.Engine
	BeforeAll(func() { router = controllers.GinEngine() })

	It("Smoke Test", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		router.ServeHTTP(w, req)
		Expect(w.Code).To(Equal(http.StatusOK))
	})

	Describe("Get user", func() {
		It("Should 400 if id is not uuid", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/user/find/1", nil)
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
		It("Should 200 if id is uuid", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/user/find/97e910d1-a572-4624-bd32-84a7a0c76bde", nil)
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))
		})
		It("Should 404 if method invalid", func() {
			cc, _ := json.Marshal(gin.H{"dsf": "sdf"})
			dd := bytes.NewBuffer(cc)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/user/find/97e910d1-a572-4624-bd32-84a7a0c76bde", dd)
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusNotFound))
		})
	})

	Describe("Register user", func() {
		It("Should able to register", func() {
			body, _ := json.Marshal(gin.H{"name": fmt.Sprint(GinkgoRandomSeed()), "password": "test"})
			buffer := bytes.NewBuffer(body)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/user/register", buffer)
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))
		})
		It("Shouldn't able to register if conflict[409]", func() {
			body, _ := json.Marshal(gin.H{"name": fmt.Sprint(GinkgoRandomSeed()), "password": "test"})
			buffer := bytes.NewBuffer(body)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/user/register", buffer)
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusConflict))
		})
		It("Shouldn't able to register if request is invalid[400]", func() {
			body, _ := json.Marshal(gin.H{"name": GinkgoRandomSeed(), "password": "test"})
			buffer := bytes.NewBuffer(body)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/user/register", buffer)
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})
})

func TestBackend(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Backend Suite")
}
