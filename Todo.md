# Environment Cleanup

`Claude.md` has been filled with planning garbage by the agent which is a very bad idea.

The agent should track changes this fork has changed in `Fork.md` and reference when needed.

`BUILD.md` and `Testing.md` have overlapping information.

`BUILD.md` is correct while `Testing.md` seems to have several halucinations such as

`--config` option being mentioned several times where filebroswer does not have this option, it does have `-c` which is correctly referenced in `BUILD.md`

Im confident `Testing.md` is nothing but confusing garbage.

# Frontend

## CSS change

I need filebrowser GUI (this repo) to match the same colour scheme used here: `https://www.acorn.tools/login`

I will use --chrome option in claude code to enable the agent to inspect the CSS for filebrowser frontend which is usually accessed with `http://localhost:8080`

You should check recent commits in git to see what has been changed so far.