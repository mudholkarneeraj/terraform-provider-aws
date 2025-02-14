package autoscaling_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccAutoScalingGroupDataSource_basic(t *testing.T) {
	datasourceName := "data.aws_autoscaling_group.test"
	resourceName := "aws_autoscaling_group.match"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, autoscaling.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "arn", resourceName, "arn"),
					resource.TestCheckResourceAttrPair(datasourceName, "availability_zones.#", resourceName, "availability_zones.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "default_cooldown", resourceName, "default_cooldown"),
					resource.TestCheckResourceAttrPair(datasourceName, "desired_capacity", resourceName, "desired_capacity"),
					resource.TestCheckResourceAttrPair(datasourceName, "health_check_grace_period", resourceName, "health_check_grace_period"),
					resource.TestCheckResourceAttrPair(datasourceName, "health_check_type", resourceName, "health_check_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "launch_configuration", resourceName, "launch_configuration"),
					resource.TestCheckResourceAttr(datasourceName, "launch_template.#", "0"),
					resource.TestCheckResourceAttrPair(datasourceName, "load_balancers.#", resourceName, "load_balancers.#"),
					resource.TestCheckResourceAttr(datasourceName, "new_instances_protected_from_scale_in", "false"),
					resource.TestCheckResourceAttrPair(datasourceName, "max_size", resourceName, "max_size"),
					resource.TestCheckResourceAttrPair(datasourceName, "min_size", resourceName, "min_size"),
					resource.TestCheckResourceAttrPair(datasourceName, "target_group_arns.#", resourceName, "target_group_arns.#"),
					resource.TestCheckResourceAttr(datasourceName, "vpc_zone_identifier", ""),
				),
			},
		},
	})
}

func TestAccAutoScalingGroupDataSource_launchTemplate(t *testing.T) {
	datasourceName := "data.aws_autoscaling_group.test"
	resourceName := "aws_autoscaling_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, autoscaling.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupDataSourceConfig_launchTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "arn", resourceName, "arn"),
					resource.TestCheckResourceAttrPair(datasourceName, "availability_zones.#", resourceName, "availability_zones.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "default_cooldown", resourceName, "default_cooldown"),
					resource.TestCheckResourceAttrPair(datasourceName, "desired_capacity", resourceName, "desired_capacity"),
					resource.TestCheckResourceAttrPair(datasourceName, "health_check_grace_period", resourceName, "health_check_grace_period"),
					resource.TestCheckResourceAttrPair(datasourceName, "health_check_type", resourceName, "health_check_type"),
					resource.TestCheckResourceAttr(datasourceName, "launch_configuration", ""),
					resource.TestCheckResourceAttrPair(datasourceName, "launch_template.#", resourceName, "launch_template.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "load_balancers.#", resourceName, "load_balancers.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "load_balancers.0.id", resourceName, "load_balancers.0.id"),
					resource.TestCheckResourceAttrPair(datasourceName, "load_balancers.0.name", resourceName, "load_balancers.0.name"),
					resource.TestCheckResourceAttrPair(datasourceName, "load_balancers.0.version", resourceName, "load_balancers.0.version"),
					resource.TestCheckResourceAttr(datasourceName, "new_instances_protected_from_scale_in", "false"),
					resource.TestCheckResourceAttrPair(datasourceName, "max_size", resourceName, "max_size"),
					resource.TestCheckResourceAttrPair(datasourceName, "min_size", resourceName, "min_size"),
					resource.TestCheckResourceAttrPair(datasourceName, "target_group_arns.#", resourceName, "target_group_arns.#"),
					resource.TestCheckResourceAttr(datasourceName, "vpc_zone_identifier", ""),
				),
			},
		},
	})
}

// Lookup based on AutoScalingGroupName
func testAccGroupDataSourceConfig_basic(rName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigLatestAmazonLinuxHvmEbsAmi(),
		acctest.ConfigAvailableAZsNoOptIn(),
		acctest.AvailableEC2InstanceTypeForAvailabilityZone("data.aws_availability_zones.available.names[0]", "t3.micro", "t2.micro"),
		fmt.Sprintf(`
data "aws_autoscaling_group" "test" {
  name = aws_autoscaling_group.match.name
}

resource "aws_autoscaling_group" "match" {
  name                      = "%[1]s_match"
  max_size                  = 0
  min_size                  = 0
  health_check_grace_period = 300
  health_check_type         = "ELB"
  desired_capacity          = 0
  force_delete              = true
  launch_configuration      = aws_launch_configuration.data_source_aws_autoscaling_group_test.name
  availability_zones        = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1]]
}

resource "aws_autoscaling_group" "no_match" {
  name                      = "%[1]s_no_match"
  max_size                  = 0
  min_size                  = 0
  health_check_grace_period = 300
  health_check_type         = "ELB"
  desired_capacity          = 0
  force_delete              = true
  launch_configuration      = aws_launch_configuration.data_source_aws_autoscaling_group_test.name
  availability_zones        = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1]]
}

resource "aws_launch_configuration" "data_source_aws_autoscaling_group_test" {
  name          = "%[1]s"
  image_id      = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = data.aws_ec2_instance_type_offering.available.instance_type
}
`, rName))
}

func testAccGroupDataSourceConfig_launchTemplate() string {
	return acctest.ConfigCompose(
		acctest.ConfigLatestAmazonLinuxHvmEbsAmi(),
		acctest.ConfigAvailableAZsNoOptIn(),
		acctest.AvailableEC2InstanceTypeForAvailabilityZone("data.aws_availability_zones.available.names[0]", "t3.micro", "t2.micro"),
		`
data "aws_autoscaling_group" "test" {
  name = aws_autoscaling_group.test.name
}

resource "aws_autoscaling_group" "test" {
  availability_zones = [data.aws_availability_zones.available.names[0]]
  desired_capacity   = 0
  max_size           = 0
  min_size           = 0
  launch_template {
    id      = aws_launch_template.test.id
    version = aws_launch_template.test.default_version
  }
}

resource "aws_launch_template" "test" {
  name_prefix   = "test"
  image_id      = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = data.aws_ec2_instance_type_offering.available.instance_type
}
`)
}
