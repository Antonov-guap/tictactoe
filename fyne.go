package main

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func drawFyneWindow(g *Game) {
	a := app.New()
	a.Settings().SetTheme(&myTheme{})
	w := a.NewWindow("Tic-Tac-Toe")
	w.Resize(fyne.NewSize(400, 367))
	w.SetFixedSize(true)
	w.CenterOnScreen()
	w.SetPadded(false)

	uri, err := storage.ParseURI("https://assetstorev1-prd-cdn.unity3d.com/key-image/7a1dfe4a-84c6-4c70-bfcb-1a81799165d3.jpg")
	logFatalOn(err)

	windowBack := canvas.NewImageFromURI(uri)
	windowBack.FillMode = canvas.ImageFillStretch

	fieldDarkBackColor := color.RGBA{R: 50, G: 50, B: 70, A: 248}
	fieldLightBackColor := color.RGBA{R: 150, G: 150, B: 170, A: 245}
	fieldBackColor := fieldLightBackColor
	if a.Settings().ThemeVariant() == theme.VariantLight {
		fieldBackColor = fieldLightBackColor
	}

	fieldBack := canvas.NewRectangle(fieldBackColor)
	fieldBack.SetMinSize(fyne.NewSize(250, 250))
	if a.Settings().ThemeVariant() == theme.VariantLight {
		fieldBack.FillColor = fieldLightBackColor
	} else {
		fieldBack.FillColor = fieldDarkBackColor
	}

	settings := make(chan fyne.Settings)
	a.Settings().AddChangeListener(settings)
	go func() {
		for s := range settings {
			if s.ThemeVariant() == theme.VariantLight {
				fieldBack.FillColor = fieldLightBackColor
			} else {
				fieldBack.FillColor = fieldDarkBackColor
			}
		}
	}()

	fieldGrid := container.NewGridWithColumns(len(g.field))

	topPanel := widget.NewLabel("Ходят: X")
	bottomPanel := widget.NewLabel("")

	var symbols = map[cell]rune{
		empty: ' ',
		cross: 'X',
		zero:  'O',
	}

	checkGameOver := func() {
		if g.IsOver() {
			for _, obj := range fieldGrid.Objects {
				button := obj.(*widget.Button)
				button.Disable()
				if g.winner == empty {
					bottomPanel.SetText("Ничья")
				} else {
					bottomPanel.SetText("Победитель - " + string(symbols[g.winner]))
				}
			}
		}
	}

	for i := range g.field {
		i := i
		for j := range g.field[i] {
			j := j
			sym := string(symbols[g.field[i][j]])
			var button *widget.Button
			button = widget.NewButton(
				sym, func() {
					err := g.MakeTurn(j, i)
					if err != nil {
						log.Printf("err: %v", err)
						return
					}
					button.SetText(string(symbols[g.field[i][j]]))
					button.Disable()
					checkGameOver()
					if g.IsOver() {
						topPanel.SetText("Игра окончена!")
					} else if g.turn == zero {
						topPanel.SetText("Ходят: O")
					} else {
						topPanel.SetText("Ходят: X")
					}
				},
			)
			fieldGrid.Add(button)
		}
	}

	field := container.NewMax(
		fieldBack,
		container.NewPadded(fieldGrid),
	)
	field.Resize(fyne.NewSize(250, 250))

	w.SetContent(
		container.NewBorder(
			container.NewCenter(topPanel),    // top
			container.NewCenter(bottomPanel), // bottom
			nil, nil,
			container.NewMax(
				windowBack,
				container.NewCenter(
					field,
				),
			),
		),
	)

	w.ShowAndRun()
}

func logFatalOn(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

type myTheme struct{}

func (a myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	th := map[fyne.ThemeVariant]map[fyne.ThemeColorName]color.Color{
		theme.VariantLight: {theme.ColorNameDisabled: color.Black},
		theme.VariantDark:  {theme.ColorNameDisabled: color.White},
	}
	if col, ok := th[variant][name]; ok {
		return col
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (a myTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (a myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (a myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
