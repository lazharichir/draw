# TODO

## Implement a tile caching mechanism
- [X] Create a function that takes a point and returns the tile it belongs to
- [X] Create a new table that stores Tile <> Changed timestamp for quick lookup
- [X] Create a function that detects changed tiles since time T
- [X] Create a function that generates updated cached tiles
- [X] Implement the TileCache
- [-] Create an endpoint that update tiles every minute

## Implement the `LandRegistry`
- [X] Define the initial interface:
  - [X] `saveLease`
  - [X] `getLeasesByPixel`
  - [X] `getLeasesByArea`

## Miscellaneaous

- [ ] Handle bulk pixels on the frontend when polling (e.g., after an image was added)
- [ ] Fix the flickering grid on the frontend