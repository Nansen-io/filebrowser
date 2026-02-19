# CSS Fix

I need filebrowser GUI (this repo) to match the same colour scheme used here: `https://www.acorn.tools/login`

I will use --chrome option in claude code to enable the agent to inspect the CSS for filebrowser frontend which is usually accessed with `http://localhost:8080`

## References

### Muted Font

`https://www.acorn.tools/login` `class="text-sm text-muted-foreground"`


### Link Font Colour

```
<a class="text-sm text-primary hover:text-primary/80 transition-colors" href="/forgot-password">Forgot Password?</a>
```

colour: `#50898e`

It is the link colour used by `https://www.acorn.tools/login`

## Current CSS Issues

See screenshot.png for examples

### File Actions

Hide

```
<button aria-label="File-Actions" class="action file-actions"><i class="material-icons">add</i> File actions</button>
```

THis is already taken care of by the actions menu on the right side bar and right click menu.

### Colours

```
<div class="inner-card"><a href="/settings#profile-main" class="person-button action button"><i class="material-icons">person</i> 16f9a1f4-4367-4d15-82b1-2e2fe44d9503 <i aria-label="settings" class="material-icons">settings</i></a></div>
```

font and icons nearly invisible against white.

use `### Link Font Colour` for link text `settings`

Icons should be same colours are File browser

Current the following are wrong:


Settings for user icons:

```
<i class="material-icons">person</i>
```

and

```
<i aria-label="settings" class="material-icons">settings</i>
```

logout icon

```
<i class="material-icons">exit_to_app</i>

```

and all icons in the right hand menu:

```
<div data-v-679ab079="" class="sidebar-actions card"><div data-v-679ab079="" class="inner-card"><button data-v-679ab079="" class="action button action-button" aria-label="New folder"><i data-v-679ab079="" class="material-icons action-icon">create_new_folder</i><span data-v-679ab079="">New folder</span></button><button data-v-679ab079="" class="action button action-button" aria-label="New file"><i data-v-679ab079="" class="material-icons action-icon">note_add</i><span data-v-679ab079="">New file</span></button><button data-v-679ab079="" class="action button action-button" aria-label="Upload"><i data-v-679ab079="" class="material-icons action-icon">file_upload</i><span data-v-679ab079="">Upload</span></button><button data-v-679ab079="" class="action button action-button" aria-label="Share"><i data-v-679ab079="" class="material-icons action-icon">share</i><span data-v-679ab079="">Share</span></button></div></div>
```