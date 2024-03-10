@echo off
setlocal
set "batch_dir=%~dp0"
echo cd %batch_dir%  & go mod init NSTU_NN_BOT & go get github.com/go-telegram-bot-api/telegram-bot-api/v5 & go get github.com/gorilla/mux
pause
endlocal