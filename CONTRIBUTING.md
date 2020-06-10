# Contribute to the RAKSH project

* [Reporting issue](#reporting-issues)
* [Pull requests](#pull-requests)
* [Golang Coding Style](#golang-coding-style)
* [Certificate of Origin](#certificate-of-origin)
* [GitHub basic setup](#github-basic-setup)
* [GitHub best practices](#github-best-practices)
* [GitHub workflow](#github-workflow)
* [Patch format](#patch-format)
* [Stable branch backports](#Stable-branch-backports)
* [Reviews](#reviews)
* [Continuous Integration](#continuous-integration)
* [Contact](#contact)
* [Project maintainers](#project-maintainers)


## Reporting Issues
Before reporting an issue, check our backlog of
[open issues](https://github.com/IBM/raksh/issues)
to see if someone else has already reported it. If so, feel free to add
your scenario, or additional information, to the discussion. Or simply
"subscribe" to it to be notified when it is updated.
If you find a new issue with the project we'd love to hear about it! The most
important aspect of a bug report is that it includes enough information for
us to reproduce it. So, please include as much detail as possible and try
to remove the extra stuff that doesn't really relate to the issue itself.
The easier it is for us to reproduce it, the faster it'll be fixed!
Please don't include any private/sensitive information in your issue!

## Pull requests

No Pull Request (PR) is too small! Typos, additional comments in the code, new testcases, bug fixes, new features, more documentation, ... it's all welcome! While bug fixes can first be identified via an "issue", that is not required. It's ok to just open up a PR with the fix, but make sure you include the same information you would have included in an issue - like how to reproduce it. PRs for new features should include some background on what use cases the new code is trying to address. When possible and when it makes sense, try to break-up larger PRs into smaller ones - it's easier to review smaller code changes. But only if those smaller ones make sense as stand-alone PRs. Regardless of the type of PR, all PRs should include:

- well documented code changes
- additional testcases. Ideally, they should fail w/o your code change applied
- documentation changes Squash your commits into logical pieces of work that might want to be reviewed separate from the rest of the PRs. But, squashing down to just one commit is ok too since in the end the entire PR will be reviewed anyway. When in doubt, squash.

PRs that fix issues should include a reference like Closes #XXXX in the commit message so that github will automatically close the referenced issue when the PR is merged.

All the repositories accept contributions via [GitHub Pull requests (PR)](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/about-pull-requests). Submit PRs by following the [GitHub workflow](#github-workflow).

## GitHub basic setup

To get started, complete the prerequisites below.

### Prerequisites

- Review [Contributor roles](#contributor-roles) that require special Git 
  configuration.

- [Set up Git](https://help.github.com/en/github/getting-started-with-github/set-up-git). 
  
  >**Note:** The email address you specify must match the email address you 
  >use to sign-off commits.

- [Fork and Clone](https://help.github.com/en/github/getting-started-with-github/fork-a-repo) the relevant repository at the 
  [Raksh Project](https://github.com/IBM/raksh).

   Example: Your local clone should show `your-github-username`, as follows.
   `https://github.com/${your-github-username}/community`.

### Contributor roles

Special Git configuration is required for these contributors:

* [Golang coding style](#golang-coding-style)

For all other contributor roles, follow the standard configuration, shown in
Prerequisites. 

### Golang coding style

* Review [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) to avoid common `Golang` errors.
* Use `gofmt` to fix any mechanical style issues.

### Dependency management
In order to add or update a dependency to this project run:
```
> export GO111MODULE=on
> go get -u [DEPENDENCY]
```
Since RAKSH uses go modules we highly recommend reading the [go modules
wiki](https://github.com/golang/go/wiki/Modules), especially the [daily workflow
section](https://github.com/golang/go/wiki/Modules#daily-workflow).
To ensure the working directory contains all necessary files afterwards, run:
```
> make vendor
```

## Certificate of Origin

In order to get a clear contribution chain of trust we use the [signed-off-by
language](https://ltsi.linuxfoundation.org/software/signed-off-process/)
used by the Linux kernel project.

## GitHub best practices

### Submit issues before PRs

Raise a GitHub issue **before** starting work on a PR.

* Our process requires an issue to be associated with every PR (see [patch format](#patch-format))

* If you are a new contributor, create an issue and add a comment stating 
that you intend to work on the issue. This notifies our team of the work you 
plan to do.

### Issue tracking

To report a bug that is not already documented, please open a GitHub issue for the repository in question.

If it is unclear which repository to raise your query against, first try to
get in [contact](#contact) with us. If in doubt, raise the issue
[here](https://github.com/IBM/raksh/issues/new) and we will
help you to handle the query by routing it to the correct area for resolution.

### Closing issues

Our tooling requires adding a `Fixes` comment to at least one commit in the PR, which triggers GitHub to automatically close the issue once the PR is merged:

```
pod: Remove token from Cmd structure

The token and pid data will be hold by the new Process structure and
they are related to a container.

Fixes #123

Signed-off-by: Ragni Garg <ragni.garg@in.ibm.com>
```

The issue is automatically closed by GitHub when the
[commit message](https://help.github.com/articles/closing-issues-via-commit-messages/) is parsed.

## GitHub workflow

Raksh employs certain augmentations to a 
[standard GitHub workflow](https://guides.github.com/introduction/flow/). 
In this section, we explain these augmentations in more detail. Follow these guidelines when contributing to Kata Containers repositories, except where noted below.

* Complete the [GitHub basic setup](#github-basic-setup) above before continuing.

* Ensure each PR only covers one topic. If you mix up different items in
  your patches or PR, they will likely need to be reworked.

* Follow a [topic branch method](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/about-branches) for development.

* Follow carefully the [patch format](#patch-format) for PRs.

* Apply the appropriate GitHub labels to your PR. This
  is particularly relevant to maintain [Stable branch backports](#stable-branch-backports). See also [GitHub labels and keywords that block PRs](#github-labels-and-keywords-that-block-prs)

  >**Note:** External contributors should use keywords, explained in the
  >link above. Labels may not be visible to external contributors. 

* [Rebase](https://help.github.com/en/github/using-git/about-git-rebase)
  commits on your branch and `force push` after each cycle of feedback.

#### Configure your environment

Most [Raksh repository](https://github.com/IBM/raksh)
contain code written in the [Go language (golang)](https://golang.org/). Go 
requires all code to be put inside the directory specified by the `$GOPATH` 
variable. Follow this example to put the code in the standard location.

```sh
$ export GOPATH=${GOPATH:-$HOME/go}
$ mkdir -p "$GOPATH"
```

>*Note*: If you intend to make minor edits, it's acceptable
> to simply fork and clone without adding the GOPATH variable.

##### Fork and clone

In this example, we configure a Git environment to contribute to this very 
`Community` repo. We create a sample branch, incorporate reviewer feedback, and rebase our commits.

1. Fork the [upstream repository](https://help.github.com/articles/cloning-a-repository):

1. [Clone your forked copy of the upstream repository](https://help.github.com/articles/cloning-a-repository):

1. While on your *forked copy*, select the green button `Clone or download`
   and copy the URL. 

1. Run the commands below and **paste the copied URL** (previous step),
   so your real GitHub user name replaces `your-github-username` below.
  
```sh
$ mkdir -p $GOPATH/src/github.com/ibm
$ cd $GOPATH/src/github.com/ibm
$ git clone https://github.com/ibm/raksh.git
$ git checkout -b 1.9.1-raksh origin/1.9.1-raksh
```
   
>**Note:** Cloning a forked repository automatically gives a remote `origin`.


##### Create a topic branch

1. Create a new "topic branch" to do your work on:

    ```
    $ git checkout -b fix-contrib-bugs
    ```

    >**Warning:** *Never* make changes directly to the `master` branch 
    >--*always* create a new "topic branch" for PR work.

1. Make some editorial changes. In this example, we modify the file that
   you are reading.

    ```
    $ $EDITOR CONTRIBUTING.md
    ```

1. Commit your changes to the current (`fix-contrib-bugs`) branch. Assure
   you use the correct [patch format](#patch-format):

    ```
    $ git commit -as
    ```

1. Push your local `fix-contrib-bugs` branch to your remote fork:

    ```
    $ git push -u origin fix-contrib-bugs
    ```

   >**Note:** The `-u` option tells `git` to "link" your local clone with
   > your remote fork so that it knows from now on that the local repository
   > and the remote fork refer to "the same" upstream repository. Strictly 
   >speaking, this option is only required the first time you call `git push`
   >for a new clone.

1. Create the PR:
  - Browse to https://github.com/IBM/raksh.
  - Click the "Compare & pull request" button that appears.
  - Click the "Create pull request" button.

  >**Note:** You do not need to change any of the defaults on this page.

##### Update your PR based on review comments

Suppose you received some reviewer feedback that asked you to make some
changes to your PR. You updated your local branch and committed those
review changes by creating three commits. There are now four commits in your
local branch: the original commit you created for the PR and three other
commits you created to apply the review changes to your branch. Your branch
now looks something like this:

```sh
$ git log master.. --oneline --decorate=no
4928d57 docs: Fix typos and fold long lines
829c6c8 apply review feedback changes
7c9b1b2 remove blank lines
60e2b2b doh - missed one
```

>**Note:** The `git log` command compares your current branch 
>(`fix-contrib-bugs`) with the `master` branch and lists all the commits,
>one per line.

##### Git rebase if multiple commits

Since all four commits are related to *the same change*, it makes sense to
combine all four commits into a *single commit* on your PR. You need to
[git rebase](https://help.github.com/github/using-git/about-git-rebase)
multiple commits on your branch. Follow these steps.

1. Update the `master` branch in your local copy of the upstream
   repository:

    ```
    $ cd $GOPATH/src/github.com/IBM/raksh
    $ git checkout master
    $ git pull --rebase upstream master
    ```

    The previous command downloads all the latest changes from the upstream
    repository and adds them to your *local copy*.

1. Now, switch back to your PR branch:

    ```
    $ git checkout fix-contrib-bugs
    ```

1. Rebase your changes against the `master` branch.

    ```sh
    $ git rebase -i master
    ```

    Example output:

    ```sh
    pick 2e335ac docs: Fix typos and fold long lines
    pick 6d6deb0 apply review feedback changes
    pick 23bc01d remove blank lines
    pick 3a4ba3f doh - missed one
    ```

1. In your editor, read the comments at the bottom of the screen. 
   Do not modify the first line, `pick 2e335ac docs: Fix typos ...`. Instead, revise `pick` to `squash` at the start of all following lines.

    Example output:

    ```sh
    pick 2e335ac docs: Fix typos and fold long lines
    squash 6d6deb0 apply review feedback changes
    squash 23bc01d remove blank lines
    squash 3a4ba3f doh - missed one
    ```

1. Save your changes and quit the editor. Git puts you *back* into your
   editor. You will see all the commit *messages*. 

1. At top is your first commit, which should be in the 
   [correct format](#patch-format). Keep your first commit and delete 
   all the following commits, as appropriate, based on 
   the review feedback.

1. Save the file and quit the editor. Once this operation completes, the four
   commits will have been converted into a single new commit. Check this by
   running the `git log` command again:

    ```sh
    $ git log master.. --oneline --decorate=no
    3ea3aef docs: Fix typo
    ```

1. Force push your updated local `fix-contrib-bugs` branch to `origin`
   remote:

    ```sh
    $ git push -f origin fix-contrib-bugs
    ```

>**Note:** Not only does this command upload your changes to your fork, it
>also includes the *latest upstream changes* to your fork since you ran
>`git pull --rebase upstream master` on the master branch and then merged
>those changes into your PR branch. This ensures your fork is now "up to
>date" with the upstream repository. The `-f` option is a "force push". Since
>you created a new commit using `git rebase`, you must "overwrite" the old
>copy of your branch in your fork on GitHub. 

Your PR is now updated on GitHub. To ensure team members are aware of this,
leave a message on the PR stating something like, "Review feedback applied".
This notification allows the team to once again review your PR more quickly.

### GitHub labels and keywords that block PRs

There are two methods that allow marking
PRs to prevent them being merged. This practice is often used during
development. The two methods are: 1) Use
[GitHub labels](https://help.github.com/articles/about-labels/)
or; 2) Use keywords in the PR subject line. The keywords can appear anywhere
in the subject line.

The following table summarises some common scenarios and appropriate use
of labels or keywords:

| Scenario | GitHub label | PR description contains |
| -------- | ------------ | ----------------------- |
| PR created "as an idea" and feedback sought | `rfc` | RFC |
| PR incomplete - needs more work or rework | `do-not-merge` `wip` | WIP |
| PR should not be merged (has all required "acks", but needs more reviewer input) | `do-not-merge` | |
| PR is a "work In progress", raised to get early feedback | `wip` | WIP |
| PR is complete but depends on another so should not be merged (yet) | `do-not-merge` | |

If any of the values in the table above are set on a PR, it will be
automatically blocked from merging.

>**Note:** Often during discussions, the abbreviated and full terms are
>used interchangeably. For instance, often `DNM` is used in discussions as 
>shorthand for `do-not-merge`. The CI systems only recognise the above 
>phrases as shown.

### Sign your PRs
The sign-off is a line at the end of the explanation for the patch. Your
signature certifies that you wrote the patch or otherwise have the right to pass
it on as an open-source patch. The rules are simple: if you can certify
the below (from [developercertificate.org](http://developercertificate.org/)):
```
Developer Certificate of Origin
Version 1.1
Copyright (C) 2004, 2006 The Linux Foundation and its contributors.
660 York Street, Suite 102,
San Francisco, CA 94110 USA
Everyone is permitted to copy and distribute verbatim copies of this
license document, but changing it is not allowed.
Developer's Certificate of Origin 1.1
By making a contribution to this project, I certify that:
(a) The contribution was created in whole or in part by me and I
    have the right to submit it under the open source license
    indicated in the file; or
(b) The contribution is based upon previous work that, to the best
    of my knowledge, is covered under an appropriate open source
    license and I have the right under that license to submit that
    work with modifications, whether created in whole or in part
    by me, under the same open source license (unless I am
    permitted to submit under a different license), as indicated
    in the file; or
(c) The contribution was provided directly to me by some other
    person who certified (a), (b) or (c) and I have not modified
    it.
(d) I understand and agree that this project and the contribution
    are public and that a record of the contribution (including all
    personal information I submit with it, including my sign-off) is
    maintained indefinitely and may be redistributed consistent with
    this project or the open source license(s) involved.
```
Then you just add a line to every git commit message:
    Signed-off-by: Joe Smith <joe.smith@email.com>
Use your real name (sorry, no pseudonyms or anonymous contributions.)
If you set your `user.name` and `user.email` git configs, you can sign your
commit automatically with `git commit -s`.

## Stable branch backports

Raksh maintains a number of stable branch releases. Bug fixes to
the master branch are selectively applied to (or "backported") these stable branches.

In order to aid identification of commits that potentially should be
backported to the stable branches, all PRs submitted must be labeled with 
one or more of the following labels. At least one label that is *not* 
`stable-candidate` must be included.

| Label              | Meaning |
| -----              | ------- |
| `bug`              | A bug fix, which will potentially be a backport candidate |
| `cleanup`          | A cleanup, which will likely not be backported                 |
| `feature`          | A new feature/enhancement, that will likely not be backported  |
| `stable-candidate` | A PR selected for backporting - very likely a bug fix          |
| `vendor`           | A golang vendor update. Might be considered for backport if the vendor update includes critical bug fixes |

In the event that a bug fix PR is selected for backporting to the stable
branches, the `stable-candidate` label is added if not already present, and
the original author of the PR is asked if they will submit the relevant
backport PRs. 

## Patch format

### General format

Beside the `Signed-off-by` footer, we expect each patch to comply with the
following format:

```
subsystem: One line change summary

More detailed explanation of your changes (why and how)
that spans as many lines as required.

A "Fixes #XXX" comment listing the GitHub issue this change resolves.
This comment is required for the main patch in a sequence. See the following examples.

Signed-off-by: <contributor@foo.com>
```

#### Pull request format

As shown above, pull requests must adhere to these guidelines:

* Preface the PR title with the appropriate keyword found in [Subsystem](#subsystem)
  
* Ensure PR title length is 75 characters or fewer, including whichever 
  `subsystem` term is used.

* Ensure the PR body line length is 72 characters or fewer.

The body of the message is **not** a continuation of the subject line and is 
not used to extend the subject line beyond its character limit. The subject 
line is a complete sentence and the body is a complete, standalone paragraph.

### Subsystem

The "subsystem" describes the area of the code that the change applies to.
It does not have to match a particular directory name in the source tree
because it is a "hint" to the reader. The subsystem is generally a single
word. Although the subsystem must be specified, it is not validated. The
author decides what is a relevant subsystem for each patch.

Examples:

| Subsystem | Description |
|--|--|
| `build` | `Makefile` or configuration script change |
| `cli` | Change affecting command line options or commands |
| `docs` | Documentation change |
| `logging` | Logging change |
| `vendor` | [Re-vendoring](#re-vendor-prs) change |

To see the subsystem values chosen for existing commits:

```
$ git log --no-merges --pretty="%s" | cut -d: -f1 | sort -u
```

### Best practices for patches

We recommend that each patch fixes one thing. Smaller patches are easier to 
review, more likely to be accepted and merged, and more conducive for
identifying problems during review.

A PR can contain multiple patches. These patches should generally be related 
to the [main patch](#main-patch) and the overall goal of the PR. However, it 
is also acceptable to include additional or 
[supplementary patches](#supplementary-patch) for things such as:

- Formatting (or whitespace) fixes
- Comment improvements
- Tidy up work
- Refactoring to simplify the codebase

### Examples

#### Main patch

The following is an example of a full patch description for the main change that shows the required "`Fixes #XXX`" comment, which references the GitHub issue this patch resolves:

```
pod: Remove token from Cmd structure

The token and pid data will be hold by the new Process structure and
they are related to a container.

Fixes: #123

Signed-off-by: Ragni Garg <ragni.garg@in.ibm.com>
```

#### Supplementary patch

If a PR contains multiple patches, [only one of those patches](#main-patch) needs to specify the "`Fixes #XXX`" comment. Supplementary patches have an identical format to the main patch, but do not need to specify a "`Fixes #XXX`"
comment.

Example:

```
image-builder: Fix incorrect error message

Fixed an error message which was referring to an incorrect rootfs
variable name.

Signed-off-by: Ragni Garg <ragni.garg@in.ibm.com>
```

## Reviews

Before your PRs are merged into the main code base, they are reviewed. We
encourage anybody to review any PR and leave feedback.

See the [PR review guide](PR-Review-Guide.md) for tips on performing a 
careful review.

We use the GitHub [Required Reviews](https://help.github.com/articles/approving-a-pull-request-with-required-reviews/)
system for reviewers to note if they agree or disagree with a PR. To have
an acknowledgment or "nack" registered with GitHub, you **must** use the
GitHub "Review changes" dialog to leave feedback. Notes left only in the
comments fields, whilst sometimes useful, will not get registered
in the acknowledgment counting system.

### Review Examples

The following is an example of a valid "ack", as long as
the "Approve" box is ticked in the Review changes dialog:

```
Excellent work - thanks for your contribution.

lgtm
```

## Continuous Integration

The Raksh project has a gating process to prevent introducing
regressions. When your PR is submitted, a Continuous Integration (CI) system
will run different checks on different platforms, based upon your changes.

Some of the checks are:

- Static analysis checks.
- Unit tests.
- Functional tests.
- Integration tests.

The Travis job will be executed right after the PR is opened, while the 
Jenkins jobs will wait to be triggered. A maintainer must add a `/test` 
comment on the PR to let the CI jobs run.

All CI jobs must pass in order to merge your PR.

## Contact

Get an [invite to our Slack channel](https://slack.k8s.io/),
  and then join us on Slack **#raksh**. 
  Feel free to reach out to our [team](https://github.com/IBM/raksh#team)

## Project maintainers

The Raksh project maintainers are the people accepting or
rejecting any PR. Although [anyone can review PRs](#reviews), only the
acknowledgement (or "ack") from an Approver counts towards the approval of a PR.

Approvers are listed in GitHub teams, one for each repository. The project
uses the
[GitHub required status checks](https://help.github.com/en/articles/enabling-required-status-checks)
along with the [GitHub `CODEOWNERS`file](https://help.github.com/en/articles/about-code-owners) to specify who can approve PRs. All repositories are configured to require:

- Two approvals from the repository-specific approval team.

#### References

- [Kata](https://github.com/kata-containers/community/blob/master/CONTRIBUTING.md)
- [Crio](https://github.com/cri-o/cri-o/blob/master/CONTRIBUTING.md)
