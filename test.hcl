


module "testmodule" {

  source = "../../mod"

  version = "1.0"

}


module "testmod1" {
  source = "https://"
}


resource "test" "test" {

  var = "${testmod1.output}"
}
