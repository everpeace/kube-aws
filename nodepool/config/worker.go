package config

type Worker struct {
	SpotFleet `yaml:"spotFleet,omitempty"`
}

type SpotFleet struct {
	TargetCapacity       int                   `yaml:"targetCapacity,omitempty"`
	SpotPrice            string                `yaml:"spotPrice,omitempty"`
	IAMFleetRoleARN      string                `yaml:"iamFleetRoleArn,omitempty"`
	LaunchSpecifications []LaunchSpecification `yaml:"launchSpecifications,omitempty"`
}

type LaunchSpecification struct {
	WeightedCapacity int    `yaml:"weightedCapacity,omitempty"`
	InstanceType     string `yaml:"instanceType,omitempty"`
	SpotPrice        string `yaml:"spotPrice,omitempty"`
	RootVolumeSize   int    `yaml:"rootVolumeSize,omitempty"`
	RootVolumeType   string `yaml:"rootVolumeType,omitempty"`
	RootVolumeIOPS   int    `yaml:"rootVolumeIOPS,omitempty"`
}

func NewDefaultWorker() Worker {
	return Worker{
		SpotFleet: newDefaultSpotFleet(),
	}
}

func newDefaultSpotFleet() SpotFleet {
	return SpotFleet{
		SpotPrice: "0.06",
		LaunchSpecifications: []LaunchSpecification{
			NewLaunchSpecification(1, "m3.medium"),
			NewLaunchSpecification(2, "m3.large"),
			NewLaunchSpecification(2, "m4.large"),
		},
	}
}

func NewLaunchSpecification(weightedCapacity int, instanceType string) LaunchSpecification {
	return LaunchSpecification{
		WeightedCapacity: weightedCapacity,
		InstanceType:     instanceType,
		RootVolumeSize:   30,
		RootVolumeIOPS:   0,
		RootVolumeType:   "gp2",
	}
}

func (f SpotFleet) Enabled() bool {
	return f.TargetCapacity > 0
}

func (f SpotFleet) IAMFleetRoleRef() string {
	if f.IAMFleetRoleARN == "" {
		return `{"Fn::Join":["", [ "arn:aws:iam::", {"Ref":"AWS::AccountId"}, ":role/aws-ec2-spot-fleet-role" ]]}`
	} else {
		return f.IAMFleetRoleARN
	}
}
