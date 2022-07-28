# Это заглушка к заданию, не успеваю сдать.
В ближайшие дни дополню commits в Pull Request.

# ДЗ 03
## п. 1 Добавить в пример с файловым сервером возможность получить список всех файлов на сервере (имя, расширение, размер в байтах)
Для получения списка файлов используем функционал `filepath.Walk` и
создаем слайс из объектов `File`.
```go
func ListDir(root string) (files []File, e error) {
	errWalk := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !info.IsDir() {
			name, ext := SplitName(info.Name())
			baseFileName := filepath.Base(path)
			fileLink := "files/" + baseFileName
			files = append(files, File{
				Link: fileLink,
				Name: name,
				Ext:  ext,
				Size: info.Size(),
			})
		}
		return nil
	})

	if errWalk != nil {
		return nil, fmt.Errorf("error walking the path %q: %v", root, errWalk)
	}

	return files, nil
}
```

Для отображения создаем отдельный handler, который по заданному шаблону выводит
полученный список файлов.

## п. 2 С помощью query-параметра, реализовать фильтрацию выводимого списка по расширению (то есть, выводить только .png файлы, или только .jpeg)
Форма для получения query-параметра реализована в шаблоне для п. 1.
И тот же handler для п. 1 позволяет фильтровать отображение файлов.
```go
func (h FilesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if files, err := fs.ListDir(h.UploadDir); err == nil {
			filter := r.URL.Query().Get("ext")
			if filter != "" {
				files = fs.FilterByExt(files, filter)
			}
			tmpl, err := template.ParseFiles(filepath.Join(h.TemplatesDir, "fileserver.html"))
			if err != nil {
				http.Error(w, "Unable to parse template", http.StatusBadRequest)
				return
			}

			if err := tmpl.Execute(w, files); err != nil {
				http.Error(w, fmt.Sprintf("cannot execute template: %s", err), http.StatusInternalServerError)
				return
			}
		}
	default:
		http.Error(w, "This method is not allowed", http.StatusMethodNotAllowed)
	}
}
```

## п. 3 *Текущая реализация сервера не позволяет хранить несколько файлов с одинаковым названием (т.к. они будут храниться в одной директории на диске). Подумайте, как  можно обойти это ограничение?
Реализована достаточно простая проверка на существование файла при его загрузке.
И если файл с заданным именем уже существует на диске, то к имени загружаемого файла добавляется цифра.
Правильным решением было бы создать отдельную структуру, сопоставляющую имя файла 
и реальный путь до него (при этом реальное имя файла на ФС должно генерироваться автоматически и быть рандомным).
И эти структуры хранить в БД, и информацию парсить уже из БД, а не из файловой системы.
Но т.к. в задании мы работаем только с файловой системой, то ограничился простым переименованием 
загружаемого файла.

## Прогресс по курсовому проекту
Пока в черновом виде проработана схема моделей для БД.