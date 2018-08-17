package handler

import (
	"net/http"
	"strings"

	"github.com/covveco/micro-learning/model"
	"github.com/covveco/micro-learning/view"
	"github.com/covveco/micro-learning/view/assets"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

type DB interface {
	GetUserByID(int) (*model.User, error)
	GetUserByToken(string) (*model.User, error)
	GetUserByStaffNo(string) (*model.User, error)
	SetUserPassword(*model.User, string) error
	CreateUser(*model.User) error
	DeleteUserByID(int) error
	GetQuestionsByCourseID(int) ([]model.Question, error)
	GetOptionsByQuestionID(int) ([]model.Option, error)
	GetCompletionsByUserID(int) ([]model.Completion, error)
	GetCompletionsByUserIDCourseID(int, int) ([]model.Completion, error)
	GetCourses() ([]model.Course, error)
	GetCourseByTemplateName(string) (*model.Course, error)
	GetCourseByID(int) (*model.Course, error)
	CreateCompletion(c *model.Completion) error
}

type Handler struct {
	DB            DB
	Log           logrus.FieldLogger
	AuthSecret    string
	AdminPassword string
}

// Routes returns a chi mux with all routes mounted.
func (h *Handler) Routes() *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(h.PrepareLog, h.LogRequest)

	// Error handling
	r.NotFound(h.NotFoundError)

	// Static files
	fileServer(r, "/public", &assetfs.AssetFS{
		Asset:     assets.Asset,
		AssetDir:  assets.AssetDir,
		AssetInfo: assets.AssetInfo,
		Prefix:    "",
	})

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(h.AuthenticateUser)
		r.Get("/", h.Root)

		r.Post("/logout", h.Logout)

		r.Get("/module", h.Course)
		r.Get("/module/{name}/{id}", h.CourseID)

		r.Get("/quiz/{id}", h.QuizID)
		r.Post("/quiz/{id}", h.QuizSubmit)

		r.Get("/certificate", h.Certificate)
		r.Get("/certificate/download", h.CertificateDownload)
		r.Get("/about", h.About)

	})

	// Admin Routes
	r.Group(func(r chi.Router) {
		r.Use(h.AuthenticateAdmin)

		r.Post("/users", h.CreateUser)
		r.Delete("/users", h.DeleteUser)
	})

	// Unprotected routes
	r.Group(func(r chi.Router) {
		r.Get("/login", h.LoginForm)
		r.Post("/login", h.Login)

		r.Get("/register/{token}", h.RegistrationForm)
		r.Post("/register", h.Registration)
	})

	return r
}

// renderView renders a view, handling a no template error with a 204 No Content
func renderView(w http.ResponseWriter, n string, d interface{}) {
	if err := view.Render(w, n, d); err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
	}
}

// redirect redirects the user.
func redirect(w http.ResponseWriter, r *http.Request, path string, message string) {
	log := r.Context().Value("log").(logrus.FieldLogger)
	log.WithFields(logrus.Fields{
		"url": path,
	}).Infof("redirecting to url")

	SetFlash(w, message)
	http.Redirect(w, r, path, http.StatusFound)
}

// fileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("fileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	redirect(w, r, "/about", "")
}
