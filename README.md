## Архиватор
### v04.09.23
---
Работает по следующим алгоритмам:
- Running
- Ngram

Для запуска нужно в качестве первого аргумента указать действие (`compress`/`decompress`), которое требуется произвести с файлом, вторым аргументом указать имя входного файла. Третьим аргументом можно указать алгоритм сжатия (по-умолчанию: `running`)

Примеры:
``` bash
archiver compress input.in # сжать файл input.in алгоритмом running, результат будет записан в файл input.in.running
archiver decompress input.in.running # восстановить файл input.in.running при помощи алгоритма running, результат будет записан в файл input.in.running.recovered
archiver compress input.in ngram # сжать файл методом ngram
archiver decompress input.in.ngram ngram # восстановить файл методом ngram 
```
