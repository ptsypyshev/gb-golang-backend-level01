# ДЗ 02
## п. 1 Добавить в приложение рассылки даты/времени возможность отправлять клиентам произвольные сообщения из консоли сервера
Предполагая, что клиентов может быть более 1, реализуем функционал обработки всех подключенных клиентов.
```golang
func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			fmt.Println(msg)
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}
```

Далее, реализуем чтение из консоли в канал.
```golang
func readMsgFromConsole(msg chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg <- scanner.Text()
	}

	if scanner.Err() != nil {
		fmt.Println("Cannot read from console!")
	}
}
```

Затем реализуем обработку соединений.
```golang
func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)
	entering <- ch
	input := bufio.NewScanner(conn)
	for input.Scan() {
	}
	leaving <- ch
	conn.Close()
}
```

И наконец функцию отправки сообщений клиенту. Если есть что-то в канале сообщений, то отправляем клиенту это сообщение.
Иначе - отправляем время как в оригинальной программе.

```golang
func clientWriter(conn net.Conn, ch <-chan string) {
	for {
		select {
		case msg := <-ch:
			_, err := io.WriteString(conn, msg)
			if err != nil {
				return
			}
		default:
			_, err := io.WriteString(conn, time.Now().Format("15:04:05\n\r"))
			if err != nil {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}
}
```

В клиент внесены небольшие доработки, очищающие буфер (иначе оставались артефакты, если текст сообщения небольшой).
```golang
func bufClear(b []byte) []byte {
	for i := 0; i < len(b); i++ {
		b[i] = ' '
	}
	return b
}
```

## п. 2 Добавить в приложение чата возможность устанавливать клиентам свой никнейм при подключении к серверу
Изменена структура приложения, основной функционал клиента и сервера вынесен в отдельные пакеты.
В пакете `chatserver` доработан метод `HandleConn`.
Добавлена возможность обрабатывать команду `/nick name`
Кроме того, при выходе клиента установлена задержка для ожидания удаления клиента из списка рассылки.
Иначе были ситуации, когда сообщение о выходе клиента отправлялось в закрытый канал.
```golang
func (s *Server) HandleConn(conn net.Conn) {
	cl := cc.NewClient(conn.RemoteAddr().String())
	go cl.Writer(conn)
	cl.MsgChan <- "You are " + cl.NickName
	cl.MsgChan <- "To set nickname use command /nick [your_nickname], for example /nick Pavel"
	s.messages <- cl.NickName + " has arrived"
	s.entering <- *cl
	input := bufio.NewScanner(conn)
	for input.Scan() {
		text := input.Text()
		if strings.Contains(text, "/nick") {
			newNickName := strings.Replace(text, "/nick ", "", 1)
			s.messages <- cl.ChangeNickName(newNickName)
		} else {
			s.messages <- fmt.Sprintf("%s: %s", cl.NickName, text)
		}
	}
	s.leaving <- *cl
	time.Sleep(time.Second) // Small pause to delete client from broadcast pool
	s.messages <- cl.NickName + " has left"
	conn.Close()
}
```

В пакете `chatclient` реализован метод `ChangeNickName`, отвечающий за смену ника у клиента.
```golang
func (cl *Client) ChangeNickName(name string) string {
	oldName := cl.NickName
	cl.NickName = name
	cl.MsgChan <- "Ok. Now you are " + cl.NickName
	return fmt.Sprintf("%s has changed its name to %s", oldName, cl.NickName)
}

```

## п. 3 *Реализовать игру “Математика на скорость”
Cервер генерирует случайное выражение с двумя операндами, сохраняет ответ, а затем отправляет выражение всем клиентам.
Первый клиент, отправивший правильный ответ - побеждает, затем генерируется следующее выражение и так далее.  

Игра реализована на базе чата из предыдущего пункта.  
Добавлен пакет `conf`, который обрабатывает аргументы запуска с помощью пакета `flag`.  
Добавлен пакет `fastmath`, в котором реализована структура с математическим заданием и ответом, а также метод генерации
случайного выражения.
```golang
func (m *MathTask) Generate() {
	rand.Seed(time.Now().UnixNano())
	a := rand.Intn(Difficulty)       //nolint:gosec // this result is not used in a secure application
	b := rand.Intn(Difficulty)       //nolint:gosec // this result is not used in a secure application
	op := rand.Intn(OperationsCount) //nolint:gosec // this result is not used in a secure application
	switch op {
	case Addition:
		m.question = fmt.Sprintf("%d + %d = ", a, b)
		m.answer = strconv.Itoa(a + b)
	case Subtraction:
		m.question = fmt.Sprintf("%d - %d = ", a, b)
		m.answer = strconv.Itoa(a - b)
	case Multiplication:
		m.question = fmt.Sprintf("%d * %d = ", a, b)
		m.answer = strconv.Itoa(a * b)
	case Division:
		m.question = fmt.Sprintf("%d / %d = ", a, b)
		m.answer = strconv.FormatFloat(float64(a)/float64(b), 'f', 1, 64) //nolint:gomnd // Standard parameters
	}
}
```

В пакете `chatserver` добавлено доп. поле в структуру `Server`, а также метод `MathTasker`,
генерирующий новый вопрос с определенной периодичностью.
```golang
func (s *Server) MathTasker() {
	for {
		if len(s.Clients) > 0 {
			s.MathTask.Generate()
			s.messages <- s.MathTask.GetQuestion()
		}
		time.Sleep(MathTaskDurationSeconds * time.Second)
	}
}
```

В метод `HandleConn` добавлено доп. условие, осуществляющее проверку ответа и выявление победителя.  
```golang
else if s.MathTask.GetQuestion() != "" && s.MathTask.GetAnswer() == text {
        s.messages <- fmt.Sprintf("Congratulations! %s answered rightly! %s%s", cl.NickName, s.MathTask.GetQuestion(), s.MathTask.GetAnswer())
        s.MathTask.SetAll("", "") // No one more player can be winner
}
```

Пакет `chatclient` остался неизменным.  
Дополнительно добавлены описания к пакетам.
