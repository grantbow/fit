#     Filesystem Issue Tracker - FIT

Filesystem Issue Tracker (FIT) conventions/format are a set of suggestions for
storing issues, one directory/folder per issue with plain text file details.

Use the editors, devices and shared file storage of your choice.

bug is a program that implements the FIT conventions and helps manage issues
versioned in git or mercurial.

The FIT format can evolve to meet changing team needs. When adapted data from
other issue systems can be imported and exported.

bug is the program written in Go developed by Dave MacFarlane (driusan).
Filesystem Issue Tracker (FIT) is the new name for the Poor Man's Issue Tracker
(PMIT) storage system also developed by driusan.

# Minimum Requirements
# Jobs
# Assumptions
# Issues
# Evolution
## Good - One File (List) For All Issues
## Better - One File Per Issue
## FIT Conventions/Format - One Directory Per Issue
### issues top directory
### issue directory name
### issue Description file
### issue tag files
### other files
### other considerations
# Tools


# Minimum Requirements

Any desktop, laptop, remote, mobile, handheld computer or phone with an editor
and access to a shared file storage can be used.

Version control using git is highly recommended to provide detailed change
history, ownership and date/times.

Other tools are not required. To speed up working with issues tools like bug
will evolve to meet existing and new needs. People should even be able to adapt
their existing tool back ends to read and write to FIT with no lost data.

# Jobs `Jobs` are the reason you use your tracker. While generally beyond the
scope of the tools used to save the issues these are what drive the assumptions
used and when building your system. How well the assumptions address your jobs
will determine how well your system is working.

# Assumptions

`Assumptions` are what we call all of the other decisions used when building
your system that are not otherwise specified. FIT conventions together with
your assumptions are intended to help naturally build systems that precisely
meet your needs. [User
Innovation](https://en.wikipedia.org/wiki/User_innovation) describes how
products and services developed and refined by end users result in better
products and services. How you system evolves is critical to the success of
it's use and of it's end users.

# Issues

The term `issue` is clear and conveys no obvious bias. We use consistent
naming for clarity. Teams have used different terms to name their issues:

    Bugs  Errors  Defects  Changes  Incidents  Problems  Requests  Tickets
    
    Troubles    Support    Calls    Cases    Tasks    Records    Alarms

    Features   Requirements   Enhancements   Events   Assignments   Q&A

    (Next) Actions       (User) Stories       (Usage) Scenarios

    Use Cases            (Backlog) Items      Agile, Lean or Kanban Cards

    Mind Map Nodes       Suggestions (Box)

We use the term `title` for the short human readable text identifier. People
may sometimes confuse the terms of `issue` and `title`. Other terms include:

    Subject    Name    Summary    Description

Records of completed issues can be very valuable. Using tag_Status_closed is
recommended but other assumptions such as moving issues to another directory
or deleting them may have other advantages or limitations.

# Evolution

Here is a brief description of how issues can be handled. This progression is
what lead to the FIT conventions.

## Good - One File For All Issues

Smaller systems are nimble but fail to meet certain needs.

To begin quickly an ISSUES.txt file or BUGS file or paper index cards
can be enough.

People love to group and evaluate issues in all sorts of ways. I would not
recommend a recursive system without clear sub-module type boundaries but lists
of issues are often a key deliverable. Issues can be arranged into groups that
may be named:

    List                 Gantt Chart          Schedule

    Project              Someday/Maybe        Follow Up

    Report               Assignments          Todo
    
    Archive              Hold                 Waiting For

    History Log          Checklist            To-Do List

    Agenda               Incubator            Weekly Review

    Kanban               Table                Calendar

    Epic                 Saga                 Mind Map


The FIT system is powerful because it has few requirements.

Sorting and listings are useful but often very context sensitive. Very
elaborate systems that include dashboards, email notifications, mobile apps,
integrations, various views, charts, filters and sorting may be useful. List
management features may also be very context sensitive.

For the specific set of issues within the context of programming code a TODO:
or BUG: comment has long been used to record issues with code. Spreading issues
across code files is awkward.

Using paper has advantages. Paper doesn't crash. Notes can be reorganized
easily. Some things are more difficult like searching, copying or keeping track
of changes. Paper is not accessible online.

Other types of electronic files often share similar attributes. "Word docs" or
spreadsheets rarely allow for enough detail and supporting materials. File
write contention can become an issue when working in teams. Coders had similar
problems before distributed version control was introduced.

## Better - One File Per Issue

One more granular step can meet more needs of small to medium sized teams.
Rules similar to FIT may be used to track one issue per file. Well formed,
leading lines of text can meet the needs of issues and tags. This fails to
separate changes to the metadata from changes to the underlying issue
description or comments.

Version control like git tracks changes very well. git versions file contents
and filesystem attributes such as directories, permissions and symbolic links.

A little more structure, just enough, makes a huge difference.

## FIT Conventions/Format - One Directory Per Issue

It is surprising how often an employer or project does not already provide an
adequate system for recording important issues in a trusted, reliable system so
that issues will not fall through the cracks.

The conventions consist of directories for each issue that contain a mandatory
Description file, tag files and any needed support files.

The few necessary parts of an issue are described for clarity.

### issues top directory

A top level `issues/` directory holds one descriptively named directory per
issue.

Human readable directory names provide context.

Directory names are short human readable `title` text identifiers. To work well
with filesystem naming conventions spaces are best replaced by dashes. Other
special characters are replaced by underscores.

Implement more than one issues directory to capture naturally similar sets of
issues. Merge and synchronize issues to adapt them as needed. Issues
naturally branch and merge along with the rest of your versioned files.

The title should never change during the life time of the issue. Tradeoffs are
involved when choosing how to store issues. Other storage naming choices may
provide better or worse tradeoffs for your needs.

Directories may contain tag files and other files as needed.

### issue directory name

Directory names are issue titles. Conventions can adapt your system to your
needs but for human use directory names should be a short human readable title
of an issue. This offers great flexibility for copying and pasting emails,
adding word docs, screen shots, videos, pdfs, voice annotations, binary files
etc.

Dashes should be interpreted as spaces and n > 1 dashes should be interpreted
as n-1 dashes when converting names to human readable text.

Underscores in names are separators. Think of them as colons without the need
to escape the file names making files easier to work with.

The interpretation of file names to text is complicated by system requirements.
Therefor no special characters (like colons : and periods . ) should be used in
file names. Again, underscores translate to separators. Dashes serve double
duty as spaces and dashes when repeated. This is a quick and easy convention
but you will not be able to place a space after a dash.

At least one "Description" file is required to contain the body of the issue.
Other supporting files and tags should be put in the directory as well.

The default title of an issue is the issue directory name.

If an issue needs to be re-titled there are several possible solutions. The
quickest is to just change the directory name. For some systems changing the
issue directory name may disrupt tools or dependencies. For some systems this
may not be a problem.

If changing the directory causes too much system disruption then any new
location needs to be well known. The new location will need to be checked every
time the title is read to override the default.

The first line of the Description file could be used to store the new title. A
convention like leading characters "title " could indicate the new title. The
rules for directory names should apply the same way.

A special tag_description_text file could be used for this purpose but this
does not seem intuitive.

The default name of the required file is "Description". A configuration option
to change the default name is in progress.

### issue Description file

Issue bodies can contain many notes, remarks, comments, ideas and
things to remember.

People like to rename things. So to prevent technical dependency problems
requiring the renaming of subdirectories the title of the issue should be
located on the first line of the "Description" file.  It will usually match the
directory name. i.e. an issue dir-name should have it's first line "title: sub
dir name"

If the first line does not begin with the letters "title" or if renaming is not
a problem for your system then the directory name is used as the title.

The rest of the file contains free form text describing the issue.  There is an
art and science to what an issue can and should contain.  It can be interpreted
as markdown format. This is the only logically required file in the directory
and the only file with a fixed file name. Contents depend upon how the system
is used but descriptions often contain the context of an issue, how to
reproduce the issue, what a desired outcome would be, etc.

Three character file extensions on Description files containing human readable
text are not recommended to make for consistent naming.

### issue id

Teams often assign a numeric id, something that will not change even if the
title does. These are usually chronological indicating an approximate issue
creation time. 

tag_id_1001

Tags are further detailed below and a feature or tool may help.

### issue tag files

Tags give this system more power and flexibility than many others.

File names beginning with 'tag' contain keys and values. Underscores in the
name separate keys and values. i.e. "tag_key_value"

The part after the first _ is the key.

The part after the second _ is the optional value. If no value is provided just
the presence of the tag_key signifies a true/false flag. Values are still
recommended for clarity.

The storage is expandable, flexible and may be updated independently of other
issue parts. Files beginning with "tag" that are not issue tag files should be
avoided in issue directories.

A special subdirectory can be used for tag keys that have implied true/present
and false/absent implied values.

Value names may optionally be stored in the file name or possibly in the file.
Implementations should trim beginning/end of line whitespace just in case.

Some example tag file names:

`tag_id_1001`
`tag_Status_done`
`tag_assigned_grantbow`
`tag_notify_driusan`
`tag_type_infrastructure`
`tag_deadline_20181019T1200`
`tag_age_stale`
`tag_help_needed`
`tag_difficulty_easy`
`tag_wtf`

Keys and values should not be case sensitive as consistent case in filesystem
names will speed up working with issues. The unix convention of all lower case
may help reduce accidental mistakes.

Three character file extensions on tag files are not recommended.

Tag files may be empty or contain arbitrary text. Tools should allow for the
preservation of the rest of the file while editing. This can be useful if a
team member includes contextual information like why this issue was given this
tag, historical comments, etc. though updating the Description text may better
meet team needs.

These tag files seem a good compromise between data processing needs and the
needs of the people using the system with or without custom tools.

Some special tags may require additional rules and/or tools. To minimize key
variances and/or value variances a tool can collect existing tag keys and
existing tag values.

Some anticipated tags might include:

`tag_Status_new`
`tag_Status_triage`
`tag_Status_confirmed`
`tag_Status_backlog`
`tag_Status_assigned`
`tag_Status_implement`
`tag_Status_fixed`
`tag_Status_peer-review`
`tag_Status_closed`

    Status values can answer the hallway question "how's that issue going?"
    People may want to know and/or may be impacted by a change made to an
    issue. Team workflows and status values vary but should have a shared
    meaning understood by team members. As status changes the git log will
    provide the date and time.

`tag_id_1001`
`tag_identifier_1001`
`tag_number_123`
`tag_fedora_12345`
`tag_debian_12345`

    *Numbering* is best done using tags. While num_1001.tag might seem
    reasonable the human readable file sorting is much easier if file names
    begin with the same prefix. The .tag extension also conflicts with DFQuery.

    While it is tempting to include numbers in issue directory names, using a
    tag allows easy renumbering. Numbers are very useful which is why many
    issue systems use them. They allow lookup and cross referencing within or
    between issue systems. Needs of small teams could require only three or
    four digits. The use of a leading one in the most significant digit can
    easily prevent confusing leading zeros.

`tag_priority_2`

    Priority values should be interpreted as numbers. Lower numbers usually
    have higher priority that naturally allow for meaningful default sort
    ordering. Priority values may change often and/or need to be automatically
    assigned.

`tag_resolution_fixed`
`tag_resolution_wont-fix`
`tag_resolution_obsolete`
`tag_resolution_duplicate`
`tag_resolution_archived`

    Outcomes are best assigned using tags. While it may be tempting to delete
    completed issues keeping them maintains important context. While version
    control can provide all deleted files as needed they are difficult to
    easily access and/or count.

`tag_stage_todo`
`tag_stage_in-progress`
`tag_stage_done`
`tag_environment_test`
`tag_environment_production`
`tag_iteration_`
`tag_milestone_`
`tag_sprint_`
`tag_week_`

    These are some tag suggestions. Build what you need. Tags may store
    information that may be better looked up in other systems but stored with
    issues for reference.

`tag_blockedby_bigger-fish-issue`
`tag_follows_previous-issue`
`tag_preceeds_next-issue`

    Tags may reference other issues and relationships between them.  Symbolic
    links are possible. The file may contain details about how these issue are
    related. Dependency tracking can be very useful in certain circumstances.

`tag_severity_`
`tag_effort_`
`tag_impact_`

    Tags may enable comparisons of relative time, work invested vs. expected
    impact. These can help evaluate team opportunity costs. Team members should
    understand the conventions used for the meanings of the values.

Tag like data may be calculated using file time stamps or other attributes.
These tags may not need to be created manually but may be derived as needed.

If a tag value is missing or invalid alternative locations can store values.
This can be implicit or explicit.

    tag_x_firstline      For positive identification the first line can be used
                         for the value of key x.

    tag_y_text           The whole file contents is used for the value of key y.


### Other Files

Human readable files accurately and quickly capture details.

### Other Considerations 

It is easy to over engineer issue trackers. Tools often must balance complexity
with ease of use within a work environment. Your choices in setting up your
system should help you meet your needs flexibly without the system feeling
unwieldy and requiring duplicate, time consuming data entry.

# Tools

Tools will likely follow to get more jobs done more easily.

Recorded issues can bring sense to an inherently disorderly process as things
move forward and/or increase the need for more features to work with issues
past and present.

The whole thing should be tracked with a version control system such as git or
another revision control systems (RCS). If using git you can call it a
distributed issue tracker. All the benefits of git and similar RCS like
synchronizing on demand allowing distributed and/or remote work, merge conflict
handling These conventions with git makes for a very capable, mobile and agile
system.

The git annex is useful for tracking larger files instead of git.

Conventions should be followed for consistent processing by a variety of tools.
Since just a few conventions are used with directories and files tools may be
written in your most convenient programming language.

A numeric running count of issues is often desired. It could be written as a
git hook to increment a tag_id_value. Quite a bit of tacit knowledge can be
quickly gained when a running identifier is assigned. Not only the total number
of issues but a feel for issue velocity can be seen by those most familiar with
the system.

Filesystems with version control provide many of the advantages that databases
provide, but when needed issues can and should be parsed into a database.
Things like querying and reporting can be more efficient using something like
Berkeley DB, Sqlite, MariaDB, Postgres, MongoDB, CouchDB, etc.

For end user access an http front end would be a good logical next step.  For
the additionally types of summaries (like BUGS.txt), reports, html files and
other files in the system exclusions can be specified so they are not checked
into the version control system.

Import and export from other issue systems could be very useful.

People, command line tools or other tool implementations (like http) should be
updated with the present conventions and assumptions in mind.
