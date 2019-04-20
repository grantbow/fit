# FAQ

Q: Why would you use this?
A: Different people may have different reasons. You can bootstrap project
   management of something you're working on without any dependencies on
   external software. Working offline has advantages. You use whatever tools you
   want (SublimeText, notepad, ed, cat, vim..) to manage issues and not be
   locked into any vendor-specific tracking software. You keep track of
   things you need to do on as granular a level as you want locally without
   management complaining about the number of tickets. You follow this
   convention and keep your issues alongside your code in the same repo. You
   can trivially define your own extensions/conventions to add to this since
   arbitrary files are permited.

Q: Why change the name from PMIT to FIT?
A: bug is the program written in Go. PMIT is the old name.
   FIT is the new name for the system that stores the issues.
   Dave MacFarlane (driusan) created bug and PMIT.

   Despite the intended meaning, I feel the word "poor" in PMIT is too
   negative. The program binary can be called by any name using using aliases.
   I also feel the word "poor" is not good description of this system.

   The "poor man's" idiom as [described](https://www.merriam-webster.com/dictionary/poor%20man's)
   is a useful description in two ways.

   First, PMIT is a less expensive system to set up and use. Ease of setup is
   an inherent benefit of FIT. Unfortunately the "poor man's" idiom used in
   PMIT is also used to indicate an entity that is similar to an original
   entity but not as talented or successful. This feels wrong to me because I
   believe using a filesystem is in some ways more talented than other systems
   and could become very successful in the future.

   I offered some patches to the original repository. Development there as
   slowed down. I have forked the original github repository to continue
   development and to improve bug which does many things well. To build the
   submodules in a golang fork the paths to the modules must be renamed.
   Therefore golang version 1.11 or greater and `export GO111MODULE=on` are
   required. This allows setting up go.mod files to rename the paths, building
   and testing. As of 2019 go.mod files are the best way to rename module
   paths.

Q: How do I add an issue x?
A: `mkdir issues/x ; edit issues/x/Description ; git add issues/x/Description ; git commit ; git push`

Q: How do I close issue x?
A: `mv issues/x/tag_status_* issues/x/tag_status_closed`

Q: How do I reopen an issue?
A: `mv issues/x/tag_status_* issues/x/tag_status_open`

Q: How do I rename an issue?
A: `git mv issues/x issues/y`

Q: How do I share issues?
A: `git clone`

Q: What if issue names conflict?
A: Why are the issue titles generic? If you're using a shared git repo to
   manage issues with other people, make sure you pull (and push) often enough
   to avoid conflicts and other users are using descriptive issue titles.

Q: What if naming issues aren't unique?
A: This is a very flexible part of the convention. Use a convention that works
   for your team. Experiment to see what works for you. Using numbers as issue
   directory names might work but may limit human readability.

Q: How do I check the issue creation time?
A: If one of the files in the issue directory has not been changed ls -l will
   show it. Checking `git log` is another way. A tag could be added for just
   this purpose.

Q: How do I check when the issue was last updated?
A: [Using find](https://stackoverflow.com/questions/5566310/how-to-recursively-find-and-list-the-latest-modified-files-in-a-directory-with-s)
   is a good way, ls -lrt | tail -1 or manually checking the maximum time in
   the issue directory. (or `git log`)

Q: How do I see the issue author?
A: `git log`

Q: How does a team comment on issues?
A: Adopt a convention to use as an extension of this standard and track it
   in a git repo. In the unlikely event that this standard becomes popular,
   whatever convention is most widely adopted/considered best practice
   should be incorporated into a future version.

   (Whatever the convention is, `git log` will almost certainly be how 
    you find the time/author of a comment.)

Q: How do I enforce [policy x] for a team
A: Write a git hook.

Q: How do I secure the system?
A: The same way you secure your files.

Q: What do you mean by a job?
A: `Jobs` are the reasons for your issues you track. This term is inspired by
   the Value Proposition Canvas by Alexander Osterwalder used with the Business
   Model Canvas. Repeatable structures will enable future tools speed up job
   completion.

   While people have other terms for "jobs" the Value Proposition Canvas talks
   about Customer Jobs vs. Products and Services. Customer Pains are relieved
   by pain relievers. Customer gains are created by Gain Creators.

   Jobs are the needs big enough to do something about or hire someone.

   There are many reasons to organize issues. The uses of a FIT storage system
   could be varied. In today's fast paced environments tasks must queue up if
   there is more to do than currently possible without more resources. Things
   can get missed if not recorded accurately.

   Computers can generate all kinds of data very quickly. Most data does not
   warrant storage in a system like this but some might.

    Help/Service/Support Desks           Software Development
    Performance        Issue Log         Bug Tracking
    Product Management Network Monitoring     System Administration
    Risk Management of Risks, Assets, Projects and Changes

   Non-computer focused areas also benefit where important communications take
   place. If you start looking it's amazing where issue type systems may be
   used.  Simple first come first served queue management like take-a-number,
   paper based type medical waiting room or food ordering systems require no
   lists or sophisticated management as the result is immediately evident. More
   sophisticated systems with many inputs and outputs can be seen through the
   lens of tracking issues. These may need to prioritize, sort, filter and view
   issues in many different ways.

    CRM         Call Center        Channel Sales Leads   Suggestion Box

    Organizational Onboarding      Customer Support

Q: Can I continuously adapt my issues?
A: Yes. Recording shared understandings can be difficult even when a system
   moves at normal speed. Teams working with important issues need trusted,
   reliable systems.

   The storage system can continuously improve, adapt and grow. Professionals
   using similar systems find that other assumptions, requirements,
   recommendations or conventions may become important after building similar
   systems, sometimes simply due to inconsistent data entry by people.
   Adjusting the system to it's inputs maintains work efficiency.

   Results and understanding problems are key drivers of success. While
   important, methods and specialization do not have the most affect on
   success. Building a system with tools that are not continuously user
   improvable stems from the common IT project disconnects between software
   tool users and tool developers. Unfortunately no one tool can do it all.
   Tools that work around different kinds of problems can not address or solve
   underlying causes. Tools need to change as problems are more clearly
   understood.

Q: What inspired FIT and what does it compare and contrast with?
A: There are some impressive options out there available as open source
   software, commercial software and SAAS cloud services providing solutions
   tracking many different kinds of issues. There are no shortage of
   solutions customized for particular purposes.

    https://www.google.com/search?q=simple+issue+tracker
    https://en.wikipedia.org/wiki/Comparison_of_issue-tracking_systems
    https://en.wikipedia.org/wiki/Comparison_of_help_desk_issue_tracking_software
    https://en.wikipedia.org/wiki/Comparison_of_CRM_systems
    https://en.wikipedia.org/wiki/Comparison_of_Mobile_CRM_systems
    https://en.wikipedia.org/wiki/Comparison_of_time_tracking_software
    https://en.wikipedia.org/wiki/Comparison_of_project_management_software
    https://dist-bugs.branchable.com/software/
    https://bugs.chromium.org/p/monorail/adminIntro

   Inspirations for these conventions, in no particular order ...

   OKR systems used by many companies including Google. It is well evangelized
   by Peter Drucker (MBO), Andy Grove of Intel, John Doerr of Kleiner Perkins
   VC firm in Menlo Park (who wrote Measure What Matters) and others.

   [David Allen's](https://davidco.com) GTD systems have inspired an amazing
   array of other useful systems. His TED talk is great. David says the art of
   stresss-free productivity flow is a martial art. Getting appropriately
   engaged is what happens during a crisis. It forces you to do the behaviors
   that could also be used without a crisis. The paradoxical truths he talks
   about can feel awkward, unnatural and unnecessary. Psychic bandwidth is key
   from the elements of control and focused attention. They come from much
   experience. His three key principles are:
    1. Get it out of your head, capture your thinking
    2. Identify and make outcome/action decisions
    3. Use the right maps

   Manager Tools guides organizational development using (un)common sense.

   Mindfulness generally as human beings can only effectively focus on a few
   things at a time.

   https://en.wikipedia.org/wiki/Lean_manufacturing

   GitLab.com builds on distributed versioning tools and has extended into a
   single application for the entire software development life cycle with
   DevOps focused features.

   Pivotal Tracker is truly amazing for high volume agile software development.

   Asana.com tracks projects well.

   Trello.com provides a great Kanban board.

   Evernote.com has an array of useful features.

   ServiceNow.com provides tools for many human workflows.

   Zenkit.com does a great job providing an elegant, easy to use system.

   Wunderlist.com is a great app for managing tasks.

   GetDoneDone.com for bugs and issues looks easy to use.

   TaskWarrior.org is a command line todo list manager.

   debbugs is still used as debian.org/Bugs starting in 1994.

   Bugzilla began use in 1998 for mozilla.org

   BestPractical.com/rt Request Tracker started as a perl and email based system in 1999.

   Bugseverywhere.org is written in python and supports many different
   distributed version control backends. It's been in use since 2005.

   RedMine.org is used in business production environments starting in 2006.

   Github and Bitbucket have hosted projects since 2008

   github.com/duplys/git-issues is very similar to bug but is written in python
   and uses a hidden branch for storage depending on a system called a shelf.
   It was able to store a version of itself under the .git/ directory so only
   python is required. I contributed some code but after that I didn't feel it
   was stable and was a bit difficult to debug because of the shelf storage.
   Development activity has decreased.

   github.com/dspinellis/git-issue is a single shell script to work with git.
   It can use an existing git repo or a new one. It stores files in a .issues
   directory as a hierarchy. It is backward compatible with gi.

   github.com/jeffWelling/ticgit is another git based system that stores data
   in a branch.

   Jonathan Corbet wrote a useful article in 2008
   https://lwn.net/Articles/281849/ about Distributed bug tracking.

   github.com/google/git-appraise is a distributed code review system written
   in go that stores reviews as git objects that reference commits using a
   built-in git notes system. https://git-scm.com/docs/git-notes

   github.com/MichaelMure/git-bug is a written in go, uses a golang struct
   (shown as json) data model to store changes to bugs as git tree (dir) and
   git blob objects in a structure like refs/bugs/<bug-id> which are hashes.
   These are aggregated into an array called an OperationPack. It uses a
   colorful, interactive terminal UI and is developing a web based UI. It has
   "bridges" to other trackers and is packaged for Archlinux.

   The issue-based information system (IBIS) is an approach developed in the
   1960's for clarifying complex, ill-defined (aka wicked) problems that
   involve multiple stakeholders. IBIS focuses on questions and is especially
   suited for exploring early phases of problem solving conversations when a
   problem is ill-defined.
   https://en.wikipedia.org/wiki/Issue-based_information_system

   [Compendium](https://en.wikipedia.org/wiki/Compendium_(software) implements
   IBIS using a notation made up of issues/questions (?), positions/ideas
   (light bulb), supporting pro arguments (+), detracting con arguments (-),
   notes (paper) and lists (bullet items). Position/ideas that show
   committment can be changed to decisions (gavel).

   Every system has strengths and weaknesses. Any of these may be of more use
   to you or more adaptable than this system under various conditions.

Q: Why do other systems use databases?
A: This is a good question. While that's just how people manage data, let's
   take a moment to step back.

   What are the advantages of databases? They are very fast and efficient if
   you know what data you want to manage. Designs and models can store large
   amounts of data. Access control is also another strength.

   However we are dealing with sometimes ad hoc data. Requirements can and do
   change. High speeds, large data sets and database managers and are not
   present or needed.

   Filesystems are one of the simplest and most pervasive ways to store data
   readily available for users. Access control is sometimes needed for
   protecting private data, but the data we are working with is mostly public
   and needs to be shared widely among team members who already have the
   necessary access restrictions in place. When combined with a version control
   system changes can be tracked and reviewed even better than with most
   database systems.

Q: How do I work with teams?
A: Members of project teams deal with seemingly unending streams of wildly
   variable `issues` generated externally or internally at any time.

   Your issue system must scale to meet a changing set of nonlinear needs.
   Personal lists of issues are great to capture issues but quickly fail to
   provide the important shared team context.

   The performance of the team is optimized by precisely, completely, quickly
   and respectfully by saving issues and team member contributions without
   detracting from keeping focus on the most impactful team goals. Recording
   ideas gives a beneficial outlet for member creativity while simultaneously
   allowing other team member views to be recorded. This helps build stronger
   shared understanding of team priorities. Wildly variable issue recordings
   can be tracked and synchronized, adding other viewpoints and perspectives.

   A set of the fewest, sufficient, shared team conventions help members work
   together to achieve shared goals. A system of shared storage is adaptable to
   many team needs and can be implemented in many ways. This system has evolved
   to capitalize on shared storage with the fewest limitations.

Q: This can never work because [...]
A: Okay, so use something that works better for you.
