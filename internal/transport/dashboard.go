package transport

import (
	"finance-tracker/internal/models"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type DashboardData struct {
	IsLoggedIn      bool
	User            models.User
	Expenses        []models.Expense
	TotalAmount     float64
	CurrentDateTime string
	FilterStart     string
	FilterEnd       string
	FilterCategory  string // Поле для сохранения выбора
	SuccessMsg      string
	ErrorMsg        string
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	filter := models.ExpenseFilter{
		UserID:    userID,
		Category:  r.URL.Query().Get("category"),
		StartDate: r.URL.Query().Get("start"),
		EndDate:   r.URL.Query().Get("end"),
	}

	user, err := h.authSvc.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	expenses, err := h.expenseSvc.GetFilteredExpenses(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := DashboardData{
		IsLoggedIn:      true,
		User:            *user,
		Expenses:        expenses,
		TotalAmount:     h.expenseSvc.CalculateTotal(expenses),
		CurrentDateTime: time.Now().Format("2006-01-02T15:04"),
		FilterStart:     filter.StartDate,
		FilterEnd:       filter.EndDate,
		FilterCategory:  filter.Category, // Передаем обратно в шаблон
		SuccessMsg:      getFlash(w, r, "success_msg"),
		ErrorMsg:        getFlash(w, r, "error_msg"),
	}

	tmpl, err := template.ParseFiles("ui/html/base.html", "ui/html/dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.ExecuteTemplate(w, "base", data)
}

func (h *Handler) AddExpense(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		setFlash(w, "error_msg", "Ошибка формы")
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	userID := r.Context().Value("userID").(int)
	category := r.FormValue("category")
	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
	desc := r.FormValue("description")
	dateStr := r.FormValue("date")

	err := h.expenseSvc.AddExpense(r.Context(), userID, category, amount, desc, dateStr)
	if err != nil {
		setFlash(w, "error_msg", "Ошибка: "+err.Error())
	} else {
		setFlash(w, "success_msg", "Расход добавлен")
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *Handler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		setFlash(w, "error_msg", "Ошибка формы")
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	userID := r.Context().Value("userID").(int)
	req := models.UpdatePasswordRequest{
		OldPassword: r.FormValue("old_password"),
		NewPassword: r.FormValue("new_password"),
	}

	err := h.authSvc.UpdatePassword(r.Context(), userID, req)
	if err != nil {
		setFlash(w, "error_msg", "Ошибка: "+err.Error())
	} else {
		setFlash(w, "success_msg", "Пароль обновлен")
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
