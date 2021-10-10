Unofficial Terraform Provider for Planetscale :rocket:
==================
- Website: https://planetscale.com
- Documentation: https://registry.terraform.io/providers/s1ntaxe770r/planetscale/latest



Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 1.0+
- Golang go1.16+

Using the provider
----------------------

See the [provider documentation](https://registry.terraform.io/providers/s1ntaxe770r/planetscale/latest) to get started
Or this [repo](https://github.com/s1ntaxe770r/tf-planetscale-db) with a full example


Contributing 
-------------
- Fork and clone this repo
- Run `make install` 
- Make and test your changes  `cd examples && make apply`
- Submit a PR!

Checkout the [roadmap](https://github.com/s1ntaxe770r/terraform-provider-planetscale/projects/2) for features i plan to add

I built this provider to answer one simple question. How hard could it be? Well turns out it's not easy but the end result here has left me with a better appreciation for the work Hashicorp has put in to Terraform, I'm still not a 100% sure how everything works but at least I'm more educated than when i started.

This provider wouldn't be possible without the amazing [client library](https://github.com/planetscale/planetscale-go) planetscale provides, i would have spent days trying to mock api requests without it.
 
If you found this cool pls leave a star for dopamine 


