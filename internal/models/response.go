package models

type ResultMessage struct {
	Tk string `json:"tk"`
	Ru string `json:"ru"`
	En string `json:"en"`
}

type DataMessage struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// todo: make correctly this response
var Success = ResultMessage{
	Tk: "Ütünlikli ýerine ýetirildi",
	Ru: "Успешно выполнено",
	En: "Successfully executed",
}

var InternalServerError = ResultMessage{
	Tk: "Serwer ýalňyşlygy",
	Ru: "Ошибка сервера",
	En: "Internal server error",
}

var InvalidInput = ResultMessage{
	Tk: "Nädogry maglumat",
	Ru: "Неверное тело запроса",
	En: "Invalid request body",
}

var ServiceUnavailableWait = ResultMessage{
	Tk: "Hayyş garaşyň",
	Ru: "пожалуйста, подождите",
	En: "please wait",
}

var UnauthorizedError = ResultMessage{
	Tk: "Hesap döredilmedik",
	Ru: "Аккаунт не создан",
	En: "Account not created",
}

var Forbitten = ResultMessage{
	Tk: "Abunaňyz ýok, yada çäkden geçdiňiz",
	Ru: "У вас нет подписки",
	En: "You don't have a subscription",
}

var Conflict = ResultMessage{
	Tk: "MAglumat eýýäm bar",
	Ru: "уже существует",
	En: "already exists",
}

var NotFound = ResultMessage{
	Tk: "Maglumat Tapylmady",
	Ru: "Ничего не найдено",
	En: "Not found",
}

var PaymentRequired = ResultMessage{
	Tk: "Toleg gerekli",
	Ru: "Требуется оплата",
	En: "Payment required",
}

type Response struct {
	Error  error       `json:"error"`
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}
