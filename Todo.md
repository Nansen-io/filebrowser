# CSS Fix

I need filebrowser GUI (this repo) to match the same colour scheme used here: `https://www.acorn.tools/login`

I will use --chrome option in claude code to enable the agent to inspect the CSS for filebrowser frontend which is usually accessed with `http://localhost:8080`

## References

### Muted Font

`https://www.acorn.tools/login` `class="text-sm text-muted-foreground"`


## Current CSS Issues

### File browser font



### 1

`class="material-icons action-icon"`

should match File browser icons:

icon forground  `#448388` 

icon background `#eaf2f4`

### 2

`class="person-button action button"`

font is to light and cannot be read, use `### Muted Font`

### 3

`class="material-icons"` which I see under `class="person-button action button"`

Icon colours are to light 

should match File browser icons:

icon forground  `#448388` 

icon background `#eaf2f4`

### 4

search bar font.

see image /home/mem/git/filebrowser/search.png

font colour is to dark, suggest `#ffffff` or something appropiate. I am unsure of this.

### 5

`class="sidebar-links card"`

should be hidden

### 6

```
<button data-v-e94d8d6f="" class="action button action-button" aria-label="New folder"><i data-v-e94d8d6f="" class="material-icons action-icon">create_new_folder</i><span data-v-e94d8d6f="">New folder</span></button>
```

Fonts in these buttons should use `### Muted Font`