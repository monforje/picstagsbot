package keyboard

import tele "gopkg.in/telebot.v4"

var MainMenu = &tele.ReplyMarkup{
	ResizeKeyboard: true,
}

var DescriptionMenu = &tele.ReplyMarkup{
	ResizeKeyboard: true,
}

var FinishUploadMenu = &tele.ReplyMarkup{
	ResizeKeyboard: true,
}

var (
	BtnUploadPhoto     = MainMenu.Text("Загрузить фото")
	BtnSearchPhoto     = MainMenu.Text("Найти фотографию")
	BtnAddDescription  = DescriptionMenu.Text("Добавить описание")
	BtnSkipDescription = DescriptionMenu.Text("Продолжить")
	BtnFinishUpload    = FinishUploadMenu.Text("Завершить")
)

func init() {
	MainMenu.Reply(
		MainMenu.Row(BtnUploadPhoto),
		MainMenu.Row(BtnSearchPhoto),
	)

	DescriptionMenu.Reply(
		DescriptionMenu.Row(BtnAddDescription, BtnSkipDescription),
	)

	FinishUploadMenu.Reply(
		FinishUploadMenu.Row(BtnFinishUpload),
	)
}
