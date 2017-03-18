# Contributing Guidelines

There are many ways that you can contribute to Miru's development, and we are always happy to receive feedback, suggestions, and contributions.

## Providing Feedback and Making Requests

The best way to have your feedback and requests heard is by starting a conversation in the **#dev-webmonitoring** channel in the [archivers.slack.com](https://archivers.slack.com/) chat [(click here for an invite)](https://archivers-slack.herokuapp.com/). Make sure to mention my username, **@zsck**. Once some discussion has been had and actionable tasks are determined, progress on tasks will be discussed in the [project issues](https://github.com/zsck/miru/issues).

## Submitting Code and Documentation

The development of new features, bug fixes, and all other improvements happens through [branches](https://help.github.com/articles/about-branches/) of Miru's `master` branch, and get merged after a [pull request](https://help.github.com/articles/about-pull-requests/) for your branch is reviewed and approved either by the project's [lead developer](https://github.com/zsck/) or at least two [other contributors](https://github.com/zsck/miru/graphs/contributors).

Code submitted for review must:

* Adhere to the [project style guide](https://github.com/zsck/miru/blob/master/docs/style-guide.md).
* Have been formatted with the [go fmt](https://golang.org/cmd/gofmt/) tool.
* Produce no errors/warnings by [go vet](https://golang.org/cmd/vet/).

We ask that contributors follow these few steps when getting started.

1. If you have not already, create a [fork](https://help.github.com/articles/fork-a-repo/) of Miru.

2. Clone your fork if you have not already.
 ```
git clone git@github.com:<your username>/miru.git
 ```

3. Create a branch based on the [issue number](https://github.com/zsck/miru/issues).

```
git checkout -b issue<issue number>
```

4. Make any changes required to resolve the issue.
5. Create your pull request, with the title being the same as the branch name.
6. Assign reviewers and work with them to make any requested changes.

After the review process is complete and all included parties agree that the code is ready to be merged, the project's [lead developer](https://github.com/zsck/) will merge your changes.