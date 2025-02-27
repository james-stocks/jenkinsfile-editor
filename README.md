# jenkinsfile-editor

## What's this?

A program for making simple programmatic changes to a [Jenkinsfile](https://www.jenkins.io/doc/book/pipeline/jenkinsfile/) without going to the lengths of having a full Groovy AST

## TODO

- [X] Read in a Jenkinsfile
- [X] Write back to text
- [ ] Make this a re-usable module
- [ ] Insert steps
- [ ] Delete steps
  - [ ] Delete surrounding stage if left empty
- [ ] Test that heredoc strings are preserved OK
- [ ] Write out with blank lines between stages
- [ ] (nice to have) preserve original whitespace, even if its untidy (i.e. don't cause a git diff beyond the intended change)