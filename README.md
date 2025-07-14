# mimismart-artylic

Проект для интеграции Artylic в Mimismart.
через httpAPI Artylic.
https://developer.arylic.com/httpapi/


# Установка
1. скачать архив релиза с github
2. распаковать архив 92.sh разместить  в /home/sh2/exe, Arylic_Ethernet.txt разчемтить в папку с скриптами
3. начать создание нового скрипта Arylic_Ethernet.txt указать IP адрес также указать имя скрипта обрабатывающий запросы к Artylic   


# Запуск
 скрипт 92.sh принимает 3тим параметром строку в hex формата ip|section|command|value 
где: ip - ip адрес Artylic устройства, section - секция команды, command - команда, value - значение команды
 пример: 192.168.20.220|player|setVolume|22 нужно перевести в hex формат 3139322e3136382e32302e3232307c706c617965727c736574566f6c756d657c3232

# Пример запуска 
```bash
./92.sh 110 22 3139322e3136382e32302e3232307c706c617965727c736574566f6c756d657c3232
