package application

import "gochat/internal/bucket"

type ChannelApplication struct {
	bucket *bucket.Bucket
}

func NewChannelApplication(bucket *bucket.Bucket) *ChannelApplication {
	return &ChannelApplication{
		bucket: bucket,
	}
}
