# Leetdoad

Fully customizable **Leet**code **d**ownl**oad**er and manager, which can download and organize your Leetcode submissions in different ways, e.g., by language, difficulty, question, etc. 

## Installation

```bash
go get github.com/jiachengxu/leetdoad
```

You can check the available flags by:

```bash
$ leetdoad -h
Usage of leetdoad:
  -config-file string
    	Path of the leetdoad config file. (default ".leetdoad.yaml")
  -cookie string
    	Cookie that used for scraping problems and solutions from Leetcode website, you can either pass it from here, or set LEETCODE_COOKIE env
  -debug
    	Debug logs.
```

Leetdoad uses a cookie to download your latest submissions from Leetcode, and cookie can be found in your browser when you visit Leetcode website. If you don't know how to find the cookie, Google is your friend.
When you have your cookie in-place, you can either pass it via `-cookie` flag or set it in the `LEETCODE_COOKIE` environment variable.

Leetdoad also requires a configuration file in YAML to organize your Leetcode submissions locally, you can either create a file, and pass it via `-config-file` flag or name it as `.leetdoad.yaml` in the folder, leetdoad will use it by default.

## Configuration

Leetdoad uses a Cookie and configuration file in YAML to scrape your Leedcode solutions, and organize them:

```yaml
# Languages defines the submissions you would like to download are written in which programming languages.
# Currently, supported values are: bash, c, cpp, golang, java, javascript, python, python3, rust, ruby, scala, swift.
# Also, it supports a special value "*", which means all above programming languages.
languages: [golang, java, cpp, python3]
# Pattern is the naming pattern for your submission file path. It has to be compatible with go template.
# Available Options:
#   - .Difficulty. The difficulty of the question.
#   - .Language. The language that the submission is written in.
#   - .QuestionName. The question name.
#   - .QuestionNumber The question number.
# Note that you don't need to specify the extension of the submission file because that will be automatically added based on the
# programming language that the submission is written in.
# For example, with the following pattern, your go and java solutions of 1.Two Sum question will be saved as:
# <current_dir>/go/1.two-sum.easy.go
# <current_dir>/java/1.two-sum.easy.go
# Note that you don't have to including all the above options in the pattern definition. 
pattern: "{{ .Language }}/{{ .QuestionNumber }}.{{ .QuestionName }}.{{ .Difficulty }}"
```

### Examples
Let's consider the following simple example, you have solved the following Leetcode questions:

| QuestionNumber | QuestionName | Difficulty | Language |
| :---------: | :---------: | :---------: | :---------: |
| 1 | Two Sum | easy | cpp, golang |
| 2 | Add Two Numbers | medium | golang, java |

#### Basic

```yaml
language: ["*"]
pattern: "{{ .QuestionNumber }}.{{ .QuestionName }}"
```

```
.
|-- 1.two-sum.cpp
|-- 1.two-sum.go
|-- 2.add-two-numbers.go
|-- 2.add-two-numbers.java
```

#### Organized by Questions

```yaml
language: ["*"]
pattern: "{{ .QuestionNumber }}.{{ .QuestionName }}/solution"
```

```
.
|-- 1.two-sum
|   |-- solution.cpp
|   |-- solution.go
|-- 2.add-two-numbers
|   |-- solution.go
|   |-- solution.java
```

#### Organized by Difficulties

```yaml
language: [golang, cpp]
pattern: "{{ .Difficulty }}/{{ .QuestionName }}"
```

```
.
|-- easy
|   |-- two-sum.cpp
|   |-- two-sum.go
|-- medium
|   |-- add-two-numbers.go
```

#### Organized by Languages

```yaml
language: [golang, java]
pattern: "{{ .Language }}/{{ .QuestionName }}-{{ .Difficulty }}"
```

```
.
|-- golang
|   |-- two-sum-easy.go
|   |-- add-two-numbers-medium.go
|-- java
|   |-- add-two-numbers-medium.java
```

## Limitations
Leetdoad scrapes submissions relatively slow. In my case, it always takes around 20 minutes to downland submissions of my ~400 solved questions. The bottleneck is not about the implementation of Leetdoad, and it is because the Leetcode API has rate-limiting configured.


