package payments

import (
	"github.com/DusanKasan/payments/internal/app/payments/internal/storage/postgres"
	"github.com/DusanKasan/payments/internal/app/payments/model"
	"github.com/gin-gonic/gin"
	errors "golang.org/x/xerrors"
	"net/http"
)

// storage represents a persistent data storage. It could/should be split in
// service and repository once the system is more complex.
type storage interface {
	// Create creates a new Payment. May return `model.ErrIdAlreadyExist` if a
	// payment with the same ID already exists in the system.
	Create(payment model.Payment) (*model.Payment, error)
	// ReadOne returns a Payment specified by its ID. It may return
	// `model.ErrIdDoesNotExist` if no Payment with such ID exists.
	ReadOne(ID string) (*model.Payment, error)
	// ReadSortedByID returns a window into the Payment collection sorted by IDs.
	// The start of this window can be moved by specifying the afterID - the
	// window will start at the first Payment with greater ID than afterID. The
	// size of the window is defined by the count parameter. It may return
	// `model.ErrIdDoesNotExist` if no Payment with such ID exists.
	ReadSortedByID(afterID *string, count int16) ([]model.Payment, error)
	// Update overwrites the Payment data for a Payment with the specified ID.
	// It may return `model.ErrIdDoesNotExist` if no Payment with such ID exists.
	Update(ID string, p model.Payment) (*model.Payment, error)
	// Delete removes the Payment with the specified ID. It may return
	// `model.ErrIdDoesNotExist` if no Payment with such ID exists.
	Delete(ID string) error
}

// handler returns a http.Handler representing the HTTP server of the payments
// sub-system. It accept storage to abstract implementation away and allow for
// internal testing. For detailed documentation of the returned server see the
// Handler function.
func handler(storage storage) http.Handler {
	r := gin.Default()

	r.GET("/payments", func(ctx *gin.Context) {
		var afterID *string
		if v, ok := ctx.GetQuery("afterID"); ok {
			afterID = &v
		}

		payments, err := storage.ReadSortedByID(afterID, 101)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		more := false
		if len(payments) == 101 {
			more = true
			payments = payments[:100]
		}

		ctx.JSON(200, gin.H{
			"payments": payments,
			"meta": gin.H{
				"more": more,
			},
		})
	})

	r.GET("/payments/:id", func(ctx *gin.Context) {
		payment, err := storage.ReadOne(ctx.Param("id"))
		if err != nil {
			if errors.Is(err, model.ErrIdDoesNotExist) {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "payment with this ID does not exist"})
				return
			}

			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, payment)
	})

	r.POST("/payments", func(ctx *gin.Context) {
		var payment model.Payment
		if err := ctx.BindJSON(&payment); err != nil {
			ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "unable to parse JSON as payment"})
			return
		}

		if err := payment.Validate(); err != nil {
			ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		created, err := storage.Create(payment)
		if err != nil {
			if errors.Is(err, model.ErrIdAlreadyExist) {
				ctx.Error(err)
				ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "payment with this ID already exist"})
				return
			}

			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusCreated, created)
	})

	r.PUT("/payments/:id", func(ctx *gin.Context) {
		var payment model.Payment
		if err := ctx.BindJSON(&payment); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "unable to parse JSON as payment"})
			return
		}

		if payment.ID != "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "payment ID must be empty in the payload"})
			return
		}
		payment.ID = ctx.Param("id")

		if err := payment.Validate(); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		updated, err := storage.Update(ctx.Param("id"), payment)
		if err != nil {
			if errors.Is(err, model.ErrIdDoesNotExist) {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "payment with this ID does not exist"})
				return
			}

			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, updated)
	})

	r.DELETE("/payments/:id", func(ctx *gin.Context) {
		if err := storage.Delete(ctx.Param("id")); err != nil {
			if errors.Is(err, model.ErrIdDoesNotExist) {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "payment with this ID does not exist"})
				return
			}

			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.Status(http.StatusOK)
	})

	return r
}

// Server returns a http.Handler that represents a HTTP server for the payments
// sub-system running with a postgres storage. It accepts the dsn connection
// string to a postgres database.
//
// The returned handler exposes the following endpoints:
// - GET /payments[?afterID=:ID]
// - GET /payments/:ID
// - POST /payments
// - PUT /payments/:ID
// - DELETE /payments/:ID
//
// For an in-depth documentation and examples see the OpenAPI definition in /api
func Handler(dsn string) (http.Handler, error) {
	s, err := postgres.New(dsn)
	if err != nil {
		return nil, errors.Errorf("unable to instantiate postgres storage: %w", err)
	}

	return handler(s), nil
}
