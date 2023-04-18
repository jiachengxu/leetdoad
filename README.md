# Leetdoad

Fully customizable **Leet**code **d**ownl**oad**er and manager, which can download and organize your Leetcode submissions in different ways, e.g., by language, difficulty, question, etc. 

- [Leetdoad](#leetdoad)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Examples](#examples)
    - [Basic](#basic)
    - [Organized by Questions](#organized-by-questions)
    - [Organized by Difficulties](#organized-by-difficulties)
    - [Organized by Languages](#organized-by-languages)
    - [GitHub Actions](#github-actions)
  - [Limitations](#limitations)

## Installation

```bash
go get github.com/jiachengxu/leetdoad
```

Alternatively, you can download the latest release from [release page](https://github.com/jiachengxu/leetdoad/releases).

You can check the available flags by:

```bash
$ leetdoad -h
Usage of leetdoad:
  -config-file string
    	Path of the leetdoad config file (default ".leetdoad.yaml")
  -cookie string
    	Cookie that used for scraping problems and solutions from Leetcode website, you can either pass it from here, or set LEETCODE_COOKIE env
  -debug
    	Debug logs
  -version
    	Show the current leetdoad version
  -header
      Include a header and footer in the scraped submission file in the style of the [VSCode LeetCode extension](https://marketplace.visualstudio.com/items?itemName=LeetCode.vscode-leetcode).
```

Leetdoad uses a cookie to download your latest submissions from Leetcode, and cookie can be found in your browser when you visit Leetcode website. If you don't know how to find the cookie, Google is your friend.
When you have your cookie in-place, you can either pass it via `-cookie` flag or set it to the `LEETCODE_COOKIE` environment variable.

Leetdoad also requires a configuration file in YAML to organize your Leetcode submissions locally, you can either create a file, and pass it via `-config-file` flag or name it as `.leetdoad.yaml` in the folder, leetdoad will use it by default.

## Configuration

```yaml
# Languages defines the submissions you would like to download are written in which programming languages.
# Currently, supported values are: bash, c, cpp, golang, java, javascript, python, python3, rust, ruby, scala, swift, typescript.
# Also, it supports a special value "*", which means all above programming languages.
languages: [golang, java, cpp, python3]
# Pattern is the naming pattern for your submission file path. It has to be compatible with go template.
# Available Options:
#   - .Difficulty - The difficulty of the question.
#   - .Language - The language that the submission is written in.
#   - .QuestionName - The question name.
#   - .QuestionNumber - The question number.
# Note that you don't need to specify the extension of the submission file because that will be automatically added based on the
# programming language that the submission is written in.
# For example, with the following pattern, your go and java solutions of 1.Two Sum question will be saved as:
# <current_dir>/go/1.two-sum.easy.go
# <current_dir>/java/1.two-sum.easy.go
# Note that you don't have to including all the above options in the pattern definition. 
pattern: "{{ .Language }}/{{ .QuestionNumber }}.{{ .QuestionName }}.{{ .Difficulty }}"
```

## Examples
If you are interested in how I organize my Leetcode submissions with Leetdoad, please check [my leetcode solutions](https://github.com/jiachengxu/oj/tree/main/leetcode), and I also use [GitHub Actions](https://github.com/jiachengxu/oj/blob/main/.github/workflows/leetdoad.yaml) to periodically fetch latest submissions from Leetcode.

Let's consider the following simple example, you have solved the following Leetcode questions:

| QuestionNumber | QuestionName | Difficulty | Language |
| :---------: | :---------: | :---------: | :---------: |
| 1 | Two Sum | easy | cpp, golang |
| 2 | Add Two Numbers | medium | golang, java |

### Basic

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

### Organized by Questions

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

### Organized by Difficulties

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

### Organized by Languages

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

### GitHub Actions
You can add the following GitHub Actions Workflows to your Leetcode submissions repository under `.github/workflows`.

```yaml
name: leetdoad
on: 
  # Trigger the workflow to update Leetcode submission every month.
  schedule:
    - cron: "0 0 1 * *"
  workflow_dispatch:
  # Trigger the workflow to update Leetcode submission on every push.
  push:
    branches:
      - main
jobs:
  leetdoad-scraping:
    runs-on: ubuntu-latest
    steps:
      # Install go in the workflow.
      - name: set up go
        uses: actions/setup-go@v2
        with:
          go-version: 1.16.1
      # Install the latest version of leetdoad.
      - name: install leetdoad
        run: go install github.com/jiachengxu/leetdoad@latest
      # Checkout to the current project root.
      - name: checkout repo
        uses: actions/checkout@v2
      # Run leetdoad command in leetcode folder.
      - name: leetdoad scraping
        run: leetdoad
        # If you want save your submissions in a subfolder of your repo, using the `working-directory` to specify relative path, otherwise remove the following line.
        working-directory: <SUB_PATH>
        env: 
          # The Cookie needs to be added in Settings -> Secrets of your repo.
          LEETCODE_COOKIE: ${{ secrets.LEETCODE_COOKIE }}
      # Commit the change.
      - name: commit change
        uses: stefanzweifel/git-auto-commit-action@v4
        with:
          commit_message: Update leetcode solutions by leetdoad.
          commit_user_email: 41898282+github-actions[bot]@users.noreply.github.com
          # Use your email address as Author.
          commit_author: Author Name <author@email.address>
```

## Limitations
Leetdoad scrapes submissions relatively slow. In my case, it always takes around 20 minutes to download submissions of my ~400 solved questions. The bottleneck is not about the implementation of Leetdoad, and it is because the Leetcode API has rate-limiting configured.


