test:
	cd $(dir) && $(MAKE) test

.PHONY: deps-install 
deps-install:
	cd $(dir) && $(MAKE) deps-install deps="$(deps)"