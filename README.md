## Releaser
A simple tool for automating release chores.

### Features
+ [Generate release notes from Pivotal Tracker](#generate-release-notes-from-pivotal-tracker)
+ [Update wiki page with release notes](#update-wiki-page-with-release-notes)


##### Generate release notes from Pivotal Tracker
First make sure `pivotal_api_token`, `pivotal_project_ids`, and `github_org` are added as top level keys in the `notes_config.json` file. Generating release notes requires the stories to be marked with a release specific label, so add on to each story in Pivotal.

```
./releaser relnotes -label rc-2018-01-16
```

##### Update wiki page with release notes
First make sure `wiki_url` is added as a top level key in the `notes_config.json` file. The `wikipage` command reads from stdin so piping in the output of the `relnotes` command works.

```
./releaser relnotes -label rc-2018-01-16 | ./releaser wp -page 2.2.13
```