.DEFAULT_GOAL = run

-include linux.mk
# -include windows.mk
-include organisations.mk

compile:
	go build -ldflags="-w -s -X 'main.CELLAR_ENDPOINT=$(CELLAR_ENDPOINT)' -X 'main.CELLAR=$(HCH01_CELLAR)' -X 'main.ODS=$(HCH01_ODS)' -X 'main.ORG_NAME=$(HCH01_NAME)'" -o launcher .

HCH01:
	go build -ldflags="-w -s -X 'main.CELLAR_ENDPOINT=$(CELLAR_ENDPOINT)' -X 'main.CELLAR=$(HCH01_CELLAR)' -X 'main.ODS=$(HCH01_ODS)' -X 'main.ORG_NAME=$(HCH01_NAME)'" -o builds/HCH01/launcher.exe .

run: compile
	./launcher