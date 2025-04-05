# Chat

## Что это
На данный момент это заготовка API для обмена сообщениями между пользователями в так называемых топиках.

## Как это использовать
Необходимо создать топик и зарегистрировать там пользователей.

Ожидается, что юзернэймы будут уникальными.

Далее на каждом пользователе можно вызвать метод `Publish` для отправки сообщений ВСЕМ пользователям топика(в том числе отправителю).

Для принятия сообщений необходимо слушать канал полученный из метода `Recieve`.

При регистрации пользователя необходимо указать `bufSize` - количество хранимых сообщений до прочтения. Если в какой-то момент количество непрочитанных сообщений превысит buf_size, то самые старые сообщения будут из очереди удаляться.

При создании топика необходимо указать `broadcasterCount` - для параллелизации доставки сообщений пользователям(полезно при большом количестве людей в топике).

## Как это работает
Цепочка отправки сообщений: 
 1. Вызов метода `Publish` на пользователе
 2. Передача сформированного пользователем сообщения в соответствующий топик
 3. Топик передает сообщение дистрибьютеру
 4. Дистрибьютор распределяет сообщение по броадкастерам
 5. Броадкастеры перенаправляют сообщение своим пользователям
 6. Сообщение лежит в буфферизованом канале, который пользователь получает из метода `Recieve`
 7. Сообщение доставлено

Простыми словами, при создании топика инициализируется дистрибюьтор, который запускает `broadcasterCount` горутин, в которых соответствующие броадкастеры слушают на выделенном им канале и перенаправляют все, что они получают в свою долю пользователей.

Дистрибьютор работает следующим образом:
 - Перенаправление сообщений тривиально
 - Для распределения пользователей по броадкастерам используется очередь наиболее свободных броадкастеров: при выходе пользователя из топика соответствующий ему броадкастер помещается в начало очереди. В очереди возможны дупликаты