# Сервис для работы с балансом пользователей

## Инструкция по запуску:

~~~
//TODO инструкция по запуску
~~~
Запросы следует отправлять на http://127.0.0.1:8080/

---
---

## Методы:

## Метод начисления/списание средств:

***POST operation/***

Производит операцию начисления или списание денежных средств на(с) счет(а) пользователя. 

### Параметры тела запроса:

Параметр|Тип|Описание|Обязательный / опциональный
---|---|---|---
id|integer|Номер пользователя.|Обязательный
type|string|Тип операции.<br>Доступные варианты: <li>'debit' - списание средств;</li><li>'credit' - зачисление средств.|Обязательный
amount|float|Сумма операции.|Обязательный

### Схема ответа:

Код статуса|Ответ
---|---
200|{"result": "success", "balance": {сумма денежных средств на счете после осуществления операции}}
400 (при недостаточности средств для проведения операции) |{"error": "not enough money"}
400 (при некорректном заполнении тела запроса)|{"error":{"amount":["amount must be positive integer more than zero"],"id":["id must be positive integer more than zero"],"type":["type must be debit or credit"]}}
400 (при отсутствии тела запроса)|{"error": "provide correct id, type and amount in request body"}
400 (при некорректном указании id пользователя|{"error": "no user with such id"}
405 (при некорректном методе запроса)|{"error": "wrong request method"}

---
---

## Метод перевода средств на счет другого пользователя:

***POST transfer/***

Переводит денежные средства со счета пользователя-отправителя на счет пользователя-получателя.

### Параметры тела запроса:

Параметр|Тип|Описание|Обязательный / опциональный
---|---|---|---
sender_id|integer|Номер пользователя-отправителя.|Обязательный
reciever_id|integer|Номер пользователя-получателя.|Обязательный
amount|float|Сумма операции.|Обязательный

### Схема ответа:

Код статуса|Ответ
---|---
200|{"result": "success", "sender_balance": {сумма денежных средств на счете отправителя после осуществления операции}, "reciever_balance": {сумма денежных средств на счете получателя после осуществления операции}}
400 (при недостаточности средств для проведения операции) |{"error": "not enough money"}
400 (при некорректном заполнении тела запроса)|{"error":{"amount":["amount must be positive integer more than zero"],"reciever_id":["reciever_id must be positive integer more than zero"],"sender_id":["sender_id must be positive integer more than zero"]}}
400 (при отсутствии тела запроса)|{"error": "provide correct sender_id, reciever_id and amount in request body"}
400 (при совпадении id отправителя и получателя{"error": "sender and receiver are the same"}
400 (при некорректном указании id пользователя|{"error": "no user with such id"}
405 (при некорректном методе запроса)|{"error": "wrong request method"}

---
---

## Метод получения баланса:

***GET balance/?id=<<user_id>>&currency=usd***

Получает баланс пользователя либо в валюте счета (RUB), либо в требуемой пользователю.

### Параметры строки запроса:

Параметр|Тип|Описание|Обязательный / опциональный
---|---|---|---
id|integer|Номер пользователя.|Обязательный
currency|string|Трехбуквенный код валюты (по умолчанию RUB)|Опциональный

### Перечень доступных валют:
AED, AFN, ALL, AMD, ANG, AOA, ARS, AUD, AWG, AZN, BAM, BBD, BDT, BGN, BHD, BIF, BMD, BND, BOB, BRL, BSD, BTC, BTN, BWP, BYN, BYR, BZD, CAD, CDF, CHF, CLF, CLP, CNY, COP, CRC, CUC, CUP, CVE, CZK, DJF, DKK, DOP, DZD, EGP, ERN, ETB, EUR, FJD, FKP, GBP, GEL, GGP, GHS, GIP, GMD, GNF, GTQ, GYD, HKD, HNL, HRK, HTG, HUF, IDR, ILS, IMP, INR, IQD, IRR, ISK, JEP, JMD, JOD, JPY, KES, KGS, KHR, KMF, KPW, KRW, KWD, KYD, KZT, LAK, LBP, LKR, LRD, LSL, LTL, LVL, LYD, MAD, MDL, MGA, MKD, MMK, MNT, MOP, MRO, MUR, MVR, MWK, MXN, MYR, MZN, NAD, NGN, NIO, NOK, NPR, NZD, OMR, PAB, PEN, PGK, PHP, PKR, PLN, PYG, QAR, RON, RSD, RWF, SAR, SBD, SCR, SDG, SEK, SGD, SHP, SLL, SOS, SRD, STD, SVC, SYP, SZL, THB, TJS, TMT, TND, TOP, TRY, TTD, TWD, TZS, UAH, UGX, USD, UYU, UZS, VEF, VND, VUV, WST, XAF, XAG, XAU, XCD, XDR, XOF, XPF, YER, ZAR, ZMK, ZMW, ZWL


### Схема ответа:

Код статуса|Ответ
---|---
200|{"id": 1, "currency": {Трехбуквенный код валюты}, "balance": {сумма денежных средств на счете получателя в запрашиваемой валюте}}
400 (при отсутствии указания id в строке запроса)|{"error": "add user id to query string"}
400 (при некорректном указании валюты)|{"error": "can't find such currency. please use another one"}
400 (при некорректном указании id пользователя|{"error": "no user with such id"}
405 (при некорректном методе запроса)|{"error": "wrong request method"}

---
---

## Метод получения истории операций:

***GET history/?id=<<user_id>>&page=1&sorted=datetime&order=asc***

Получения списка операций пользователя. На странице выводятся 10 операций.

### Параметры строки запроса:

Параметр|Тип|Описание|Обязательный / опциональный
---|---|---|---
id|integer|Номер пользователя.|Обязательный
page|integer|Номер страницы (по умолчанию 1)|Опциональный
sorted|string|Поле сортировки (по умолчанию datetime)<br>Доступные варианты: <li>'datetime' - время операции;</li><li>'amount' - сумма операции.|Опциональный
order|string|Порядок сортировки (по умолчанию ASC)<br>Доступные варианты: <li>'ASC' - по возрастанию;</li><li>'DESC' - по убыванию.|Опциональный

### Схема ответа:

Код статуса|Ответ
---|---
200|{"user_id": 2, "type": "debit", "amount": 10, "datetime": "2022-03-13T16:11:39Z", "idx": 11}
400 (при отсутствии указания id в строке запроса)|{"error": "add user id to query string"}
400 (при некорректном указании порядка сортировки)|{"error": "order must be asc or desc"}
400 (при некорректном указании поля сортировки)|{"error": "sorted must be amount or datetime"}
400 (при некорректном указании id пользователя|{"error": "no user with such id"}
404 (при отсутствии истории операций)|{"error": ""nothing found""}
405 (при некорректном методе запроса)|{"error": "wrong request method"}

//TODO sql-скрипт создания необходимых таблиц