# CSS Fix

I need filebrowser GUI (this repo) to match the same colour scheme used here: `https://www.acorn.tools/login`

I will use --chrome option in claude code to enable the agent to inspect the CSS for filebrowser frontend which is usually accessed with `http://localhost:8080`

## Current CSS Issues

### Icons

```
<div class="quick-toggles"><div class="clickable"><i class="material-icons">ads_click</i></div><div aria-label="Toggle Theme" class="clickable"><i class="material-icons">dark_mode</i></div><div class="clickable active"><i class="material-icons">push_pin</i></div></div>
```

class `clickable` when not enabled / toggled on the background colour is to light.

it should be `#eaf2f5` (selected using eyedropper plugin)

### Sidebar Actions Menu

should have same CSS style as 

```
<div data-v-d98731eb="" id="context-menu" 
```

currently until user hovers over it is all white.

### Hamburger Menu broken

```
<svg id="button-toggle-navbar" class="ham hamRotate180 ham5" viewBox="0 0 100 100" width="30"><path class="line top" d="M 30,33 H 70"></path><path class="line middle" d="M 30,50 H 70"></path><path class="line bottom" d="M 30,67 H 70"></path></svg>
```

Hamburger menu should be disabled and hidden.

### Search bar

```
id="search" 
```

is not visible until user clicks on it, CSS mismatch.

search dialog popup also has CSS colours mismatch