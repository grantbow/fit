
**TOC:**

<!-- toc -->

- [Name](#name)
- [History](#history)
- [Background](#Background)

<!-- tocstop -->

## Name

"fit" can mean anything, depending on your mood:

* random three letter, pronouncable characters not used by other common UNIX commands.
* filesystem issue tracker
* f\* ignorant transgression when it breaks
* strong system with stamina, flexibility and endurance

## History

fit is the golang program first developed as "bug" by Dave MacFarlane (driusan).
Filesystem Issue Tracker ([Filesystem_Issues.md](docs/Filesystem_Issues.md)) is the new name for the Poor Man's
Issue Tracker (PMIT) storage system also first developed by driusan. See the
2016 demo video of [driusan's
talk](https://www.youtube.com/watch?v=ysgMlGHtDMo) at the first
GolangMontreal.org conference, GoMTL-01. The program and storage system have incrementally
evolved while trying to remain backward compatible. See the docs/[FAQ.md](docs/FAQ.md)
for even more information.

## Background

A limited but sufficient number of conventions with just enough organization
can quickly capture issues using human readable issue directories and files.
fit can be the primary system if no other system is provided or supplement
other issue/bug systems to quickly capture issues and their context as close
to the code as possible.

Using fit helps implementers streamline working with
[issues](https://en.wikipedia.org/wiki/Issue_tracking_system) and [version
control](https://en.wikipedia.org/wiki/Version_control). fit works with
both git and mercurial distributed version control though the git features are
more well exercised.

fit is designed to adapt to your processes using issue key/value pair metadata.

The fit tool manages issues using conventions/format of
Filesystem Issue Tracker (see [Filesystem_Issues.md](docs/Filesystem_Issues.md)). A `fit/` or `issues/`
directory holds one descriptively titled directory per issue. Each directory 
holds a file Description (name is configurable) which is a text file.
Issue directories hold anything else needed about the issue.

Issue systems typically evolve from the most simple systems that work to slightly
more complex systems that work better when working with others.

At first people may naturally try to keep track of issues in a single text
file and/or spreadsheet but these can fail to meet project needs.
(see docs/[FAQ.md](docs/FAQ.md))

Issue context is valuable to coders and may be difficult for others to
understand, especially without the context of the code they describe.
fit can support multiple `fit/` directories in a
repository's tree for stronger coordination of coding and issue tracking.

Projects in IT environments face all too common circumstances: implementers
may not be given the tools needed (or given bad tools) to record code issues.
Other issue systems typically take some resources to setup and maintain which
can be difficult to justify until long after a system is desperately needed
by implementers. Separate issue systems are often focused on higher volume
user facing streams of problem reports. These systems may or may not be meet
project needs to capture valuable implementation details so valuable details
are often poorly documented or completely lost. Some IT groups hope that
problems will not attract attention and see obfuscation as a way to reduce
perhaps already oversized workloads. It can be reasonable to focus valuable
project budget, time, scope, quality or other resources on new features but
code level or low volume streams of issues can be ignored. Beyond managing
contentious flat files or spreadsheets there is FIT.

Important issues can be captured, surfaced and addressed, whether they are
actual problems, questions, possible features or ideas by those most familiar
with the project. It is hoped that all code savvy project collaborators can
capture implementation details of varying importance quickly and easily using
fit compared to using larger, possibly distracting systems best designed for
other uses like code reviews or operational facing streams of issues. Regardless
of other available issue systems that might be available, using a fit system
might advantageously complement project workflows.

[Software Development Life Cycles](https://en.wikipedia.org/wiki/Software_development_process) (SDLCs) involve more than just the source code.
Over time needs of a project may change from hacking/coding, just getting something working,
to implementing more disciplined software engineering best practices. Code can
start small and grow as users, TODO comments, use cases and developers are added.
The FIT issue system was designed to adapt to each stage of needs.

While one issue set used for one git repository may be enough the use of
recursive fit directories are now supported. As complexity increases adding
multiple `fit/` directories in different parts of your git repo may help
project coders keep focused.

There are some choices each project can make for how to handle closed
issues. As the number of issues grows closed issues can simply be deleted or
an archive can hold the inactive issues. While deleting issues helps keep things
uncluttered issues still have value over time and may be difficult to find
using only version control history.

fit can be aliased as a git subcommand "git fit ..." It is intended that similar
subcommands perform similarly expected functions.

fit software is written using [golang](https://golang.org) to make things easy,
simple and reliable. Why go? [This video](https://vimeo.com/69237265) from a
2013 Ruby conference by Andrew Gerrand of Google seems a good explanation. It
was linked from the golang.org home page.

Engineers know that there is more to code than the source itself. For some rare
individuals the code is enough context. For most people new to a code base or
distracted by other concerns any recorded context can be extremely helpful.
Notes about refactoring history, code reviews or feature ideas can be
important to grok a code base more quickly. This context may originate from
researching a user reported problem or may arise any time while coding.

