# TODO

## Implement a tile caching mechanism
- [ ] Create a function that detects changed tiles since time T
- [ ] Create a function that generates updated cached tiles
- [ ] Create an endpoint that update tiles every minute

## Implement the `LandRegistry`
- [ ] Define the initial `Store` interface:
  - [ ] `saveLease`
  - [ ] `getLeasesByPixel`
  - [ ] `getLeasesByArea`

## Miscellaneaous

- [ ] Handle bulk pixels on the frontend when polling (e.g., after an image was added)
- [ ] Fix the flickering grid on the frontend