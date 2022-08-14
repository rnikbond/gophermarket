#!/bin/bash

# Запуск системы расчета баллов лояльности
./accrual_linux_amd64 -a=:8000 &

# Запуск сервера
./gophermarket
