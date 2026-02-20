# Goal

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


## Issues

### Hide

In right side bar menu

```
<div data-v-3237c3ef="" class="sidebar-actions card">
```


hide

```
<button data-v-3237c3ef="" class="action button action-button" aria-label="Share"><i data-v-3237c3ef="" class="material-icons action-icon">share</i><span data-v-3237c3ef="">Share</span></button>
```

### Colours


```
<a href="/settings#profile-main" class="person-button action button"><i class="material-icons">person</i> 16f9a1f4-4367-4d15-82b1-2e2fe44d9503 <i aria-label="settings" class="material-icons">settings</i></a>
```

The text `16f9a1f4-4367-4d15-82b1-2e2fe44d9503` is the wrong colour.

use `### Muted Font` 