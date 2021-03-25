Generated from https://github.com/Azure/azure-rest-api-specs/tree/92ab22b49bd085116af0c61fada2c6c360702e9e/specification/mediaservices/resource-manager/readme.md tag: `package-2020-05`

Code generator @microsoft.azure/autorest.go@2.1.175


## Breaking Changes

### Removed Constants

1. Priority.Normal
1. StorageAuthentication.System

## Signature Changes

### Const Types

1. High changed type from Priority to BlurType
1. Low changed type from Priority to BlurType
1. ManagedIdentity changed type from StorageAuthentication to CreatedByType

### New Constants

1. AttributeFilter.All
1. AttributeFilter.Bottom
1. AttributeFilter.Top
1. AttributeFilter.ValueEquals
1. BlurType.Black
1. BlurType.Box
1. BlurType.Med
1. ChannelMapping.BackLeft
1. ChannelMapping.BackRight
1. ChannelMapping.Center
1. ChannelMapping.FrontLeft
1. ChannelMapping.FrontRight
1. ChannelMapping.LowFrequencyEffects
1. ChannelMapping.StereoLeft
1. ChannelMapping.StereoRight
1. CreatedByType.Application
1. CreatedByType.Key
1. CreatedByType.User
1. EncoderNamedPreset.H265AdaptiveStreaming
1. EncoderNamedPreset.H265ContentAwareEncoding
1. EncoderNamedPreset.H265SingleBitrate1080p
1. EncoderNamedPreset.H265SingleBitrate4K
1. EncoderNamedPreset.H265SingleBitrate720p
1. FaceRedactorMode.Analyze
1. FaceRedactorMode.Combined
1. FaceRedactorMode.Redact
1. H265Complexity.H265ComplexityBalanced
1. H265Complexity.H265ComplexityQuality
1. H265Complexity.H265ComplexitySpeed
1. H265VideoProfile.H265VideoProfileAuto
1. H265VideoProfile.H265VideoProfileMain
1. OdataTypeBasicCodec.OdataTypeMicrosoftMediaH265Video
1. OdataTypeBasicInputDefinition.OdataTypeInputDefinition
1. OdataTypeBasicInputDefinition.OdataTypeMicrosoftMediaFromAllInputFile
1. OdataTypeBasicInputDefinition.OdataTypeMicrosoftMediaFromEachInputFile
1. OdataTypeBasicInputDefinition.OdataTypeMicrosoftMediaInputFile
1. OdataTypeBasicJobInput.OdataTypeMicrosoftMediaJobInputSequence
1. OdataTypeBasicLayer.OdataTypeMicrosoftMediaH265Layer
1. OdataTypeBasicLayer.OdataTypeMicrosoftMediaH265VideoLayer
1. OdataTypeBasicTrackDescriptor.OdataTypeMicrosoftMediaAudioTrackDescriptor
1. OdataTypeBasicTrackDescriptor.OdataTypeMicrosoftMediaSelectAudioTrackByAttribute
1. OdataTypeBasicTrackDescriptor.OdataTypeMicrosoftMediaSelectAudioTrackByID
1. OdataTypeBasicTrackDescriptor.OdataTypeMicrosoftMediaSelectVideoTrackByAttribute
1. OdataTypeBasicTrackDescriptor.OdataTypeMicrosoftMediaSelectVideoTrackByID
1. OdataTypeBasicTrackDescriptor.OdataTypeMicrosoftMediaVideoTrackDescriptor
1. OdataTypeBasicTrackDescriptor.OdataTypeTrackDescriptor
1. Priority.PriorityHigh
1. Priority.PriorityLow
1. Priority.PriorityNormal
1. StorageAuthentication.StorageAuthenticationManagedIdentity
1. StorageAuthentication.StorageAuthenticationSystem
1. TrackAttribute.Bitrate
1. TrackAttribute.Language

### New Funcs

1. *FromAllInputFile.UnmarshalJSON([]byte) error
1. *FromEachInputFile.UnmarshalJSON([]byte) error
1. *InputDefinition.UnmarshalJSON([]byte) error
1. *InputFile.UnmarshalJSON([]byte) error
1. *JobInputSequence.UnmarshalJSON([]byte) error
1. AacAudio.AsH265Video() (*H265Video, bool)
1. Audio.AsH265Video() (*H265Video, bool)
1. AudioTrackDescriptor.AsAudioTrackDescriptor() (*AudioTrackDescriptor, bool)
1. AudioTrackDescriptor.AsBasicAudioTrackDescriptor() (BasicAudioTrackDescriptor, bool)
1. AudioTrackDescriptor.AsBasicTrackDescriptor() (BasicTrackDescriptor, bool)
1. AudioTrackDescriptor.AsBasicVideoTrackDescriptor() (BasicVideoTrackDescriptor, bool)
1. AudioTrackDescriptor.AsSelectAudioTrackByAttribute() (*SelectAudioTrackByAttribute, bool)
1. AudioTrackDescriptor.AsSelectAudioTrackByID() (*SelectAudioTrackByID, bool)
1. AudioTrackDescriptor.AsSelectVideoTrackByAttribute() (*SelectVideoTrackByAttribute, bool)
1. AudioTrackDescriptor.AsSelectVideoTrackByID() (*SelectVideoTrackByID, bool)
1. AudioTrackDescriptor.AsTrackDescriptor() (*TrackDescriptor, bool)
1. AudioTrackDescriptor.AsVideoTrackDescriptor() (*VideoTrackDescriptor, bool)
1. AudioTrackDescriptor.MarshalJSON() ([]byte, error)
1. Codec.AsH265Video() (*H265Video, bool)
1. CopyAudio.AsH265Video() (*H265Video, bool)
1. CopyVideo.AsH265Video() (*H265Video, bool)
1. FromAllInputFile.AsBasicInputDefinition() (BasicInputDefinition, bool)
1. FromAllInputFile.AsFromAllInputFile() (*FromAllInputFile, bool)
1. FromAllInputFile.AsFromEachInputFile() (*FromEachInputFile, bool)
1. FromAllInputFile.AsInputDefinition() (*InputDefinition, bool)
1. FromAllInputFile.AsInputFile() (*InputFile, bool)
1. FromAllInputFile.MarshalJSON() ([]byte, error)
1. FromEachInputFile.AsBasicInputDefinition() (BasicInputDefinition, bool)
1. FromEachInputFile.AsFromAllInputFile() (*FromAllInputFile, bool)
1. FromEachInputFile.AsFromEachInputFile() (*FromEachInputFile, bool)
1. FromEachInputFile.AsInputDefinition() (*InputDefinition, bool)
1. FromEachInputFile.AsInputFile() (*InputFile, bool)
1. FromEachInputFile.MarshalJSON() ([]byte, error)
1. H264Layer.AsBasicH265VideoLayer() (BasicH265VideoLayer, bool)
1. H264Layer.AsH265Layer() (*H265Layer, bool)
1. H264Layer.AsH265VideoLayer() (*H265VideoLayer, bool)
1. H264Video.AsH265Video() (*H265Video, bool)
1. H265Layer.AsBasicH265VideoLayer() (BasicH265VideoLayer, bool)
1. H265Layer.AsBasicLayer() (BasicLayer, bool)
1. H265Layer.AsBasicVideoLayer() (BasicVideoLayer, bool)
1. H265Layer.AsH264Layer() (*H264Layer, bool)
1. H265Layer.AsH265Layer() (*H265Layer, bool)
1. H265Layer.AsH265VideoLayer() (*H265VideoLayer, bool)
1. H265Layer.AsJpgLayer() (*JpgLayer, bool)
1. H265Layer.AsLayer() (*Layer, bool)
1. H265Layer.AsPngLayer() (*PngLayer, bool)
1. H265Layer.AsVideoLayer() (*VideoLayer, bool)
1. H265Layer.MarshalJSON() ([]byte, error)
1. H265Video.AsAacAudio() (*AacAudio, bool)
1. H265Video.AsAudio() (*Audio, bool)
1. H265Video.AsBasicAudio() (BasicAudio, bool)
1. H265Video.AsBasicCodec() (BasicCodec, bool)
1. H265Video.AsBasicImage() (BasicImage, bool)
1. H265Video.AsBasicVideo() (BasicVideo, bool)
1. H265Video.AsCodec() (*Codec, bool)
1. H265Video.AsCopyAudio() (*CopyAudio, bool)
1. H265Video.AsCopyVideo() (*CopyVideo, bool)
1. H265Video.AsH264Video() (*H264Video, bool)
1. H265Video.AsH265Video() (*H265Video, bool)
1. H265Video.AsImage() (*Image, bool)
1. H265Video.AsJpgImage() (*JpgImage, bool)
1. H265Video.AsPngImage() (*PngImage, bool)
1. H265Video.AsVideo() (*Video, bool)
1. H265Video.MarshalJSON() ([]byte, error)
1. H265VideoLayer.AsBasicH265VideoLayer() (BasicH265VideoLayer, bool)
1. H265VideoLayer.AsBasicLayer() (BasicLayer, bool)
1. H265VideoLayer.AsBasicVideoLayer() (BasicVideoLayer, bool)
1. H265VideoLayer.AsH264Layer() (*H264Layer, bool)
1. H265VideoLayer.AsH265Layer() (*H265Layer, bool)
1. H265VideoLayer.AsH265VideoLayer() (*H265VideoLayer, bool)
1. H265VideoLayer.AsJpgLayer() (*JpgLayer, bool)
1. H265VideoLayer.AsLayer() (*Layer, bool)
1. H265VideoLayer.AsPngLayer() (*PngLayer, bool)
1. H265VideoLayer.AsVideoLayer() (*VideoLayer, bool)
1. H265VideoLayer.MarshalJSON() ([]byte, error)
1. Image.AsH265Video() (*H265Video, bool)
1. InputDefinition.AsBasicInputDefinition() (BasicInputDefinition, bool)
1. InputDefinition.AsFromAllInputFile() (*FromAllInputFile, bool)
1. InputDefinition.AsFromEachInputFile() (*FromEachInputFile, bool)
1. InputDefinition.AsInputDefinition() (*InputDefinition, bool)
1. InputDefinition.AsInputFile() (*InputFile, bool)
1. InputDefinition.MarshalJSON() ([]byte, error)
1. InputFile.AsBasicInputDefinition() (BasicInputDefinition, bool)
1. InputFile.AsFromAllInputFile() (*FromAllInputFile, bool)
1. InputFile.AsFromEachInputFile() (*FromEachInputFile, bool)
1. InputFile.AsInputDefinition() (*InputDefinition, bool)
1. InputFile.AsInputFile() (*InputFile, bool)
1. InputFile.MarshalJSON() ([]byte, error)
1. JobInput.AsJobInputSequence() (*JobInputSequence, bool)
1. JobInputAsset.AsJobInputSequence() (*JobInputSequence, bool)
1. JobInputClip.AsJobInputSequence() (*JobInputSequence, bool)
1. JobInputHTTP.AsJobInputSequence() (*JobInputSequence, bool)
1. JobInputSequence.AsBasicJobInput() (BasicJobInput, bool)
1. JobInputSequence.AsBasicJobInputClip() (BasicJobInputClip, bool)
1. JobInputSequence.AsJobInput() (*JobInput, bool)
1. JobInputSequence.AsJobInputAsset() (*JobInputAsset, bool)
1. JobInputSequence.AsJobInputClip() (*JobInputClip, bool)
1. JobInputSequence.AsJobInputHTTP() (*JobInputHTTP, bool)
1. JobInputSequence.AsJobInputSequence() (*JobInputSequence, bool)
1. JobInputSequence.AsJobInputs() (*JobInputs, bool)
1. JobInputSequence.MarshalJSON() ([]byte, error)
1. JobInputs.AsJobInputSequence() (*JobInputSequence, bool)
1. JpgImage.AsH265Video() (*H265Video, bool)
1. JpgLayer.AsBasicH265VideoLayer() (BasicH265VideoLayer, bool)
1. JpgLayer.AsH265Layer() (*H265Layer, bool)
1. JpgLayer.AsH265VideoLayer() (*H265VideoLayer, bool)
1. Layer.AsBasicH265VideoLayer() (BasicH265VideoLayer, bool)
1. Layer.AsH265Layer() (*H265Layer, bool)
1. Layer.AsH265VideoLayer() (*H265VideoLayer, bool)
1. PngImage.AsH265Video() (*H265Video, bool)
1. PngLayer.AsBasicH265VideoLayer() (BasicH265VideoLayer, bool)
1. PngLayer.AsH265Layer() (*H265Layer, bool)
1. PngLayer.AsH265VideoLayer() (*H265VideoLayer, bool)
1. PossibleAttributeFilterValues() []AttributeFilter
1. PossibleBlurTypeValues() []BlurType
1. PossibleChannelMappingValues() []ChannelMapping
1. PossibleCreatedByTypeValues() []CreatedByType
1. PossibleFaceRedactorModeValues() []FaceRedactorMode
1. PossibleH265ComplexityValues() []H265Complexity
1. PossibleH265VideoProfileValues() []H265VideoProfile
1. PossibleOdataTypeBasicInputDefinitionValues() []OdataTypeBasicInputDefinition
1. PossibleOdataTypeBasicTrackDescriptorValues() []OdataTypeBasicTrackDescriptor
1. PossibleTrackAttributeValues() []TrackAttribute
1. SelectAudioTrackByAttribute.AsAudioTrackDescriptor() (*AudioTrackDescriptor, bool)
1. SelectAudioTrackByAttribute.AsBasicAudioTrackDescriptor() (BasicAudioTrackDescriptor, bool)
1. SelectAudioTrackByAttribute.AsBasicTrackDescriptor() (BasicTrackDescriptor, bool)
1. SelectAudioTrackByAttribute.AsBasicVideoTrackDescriptor() (BasicVideoTrackDescriptor, bool)
1. SelectAudioTrackByAttribute.AsSelectAudioTrackByAttribute() (*SelectAudioTrackByAttribute, bool)
1. SelectAudioTrackByAttribute.AsSelectAudioTrackByID() (*SelectAudioTrackByID, bool)
1. SelectAudioTrackByAttribute.AsSelectVideoTrackByAttribute() (*SelectVideoTrackByAttribute, bool)
1. SelectAudioTrackByAttribute.AsSelectVideoTrackByID() (*SelectVideoTrackByID, bool)
1. SelectAudioTrackByAttribute.AsTrackDescriptor() (*TrackDescriptor, bool)
1. SelectAudioTrackByAttribute.AsVideoTrackDescriptor() (*VideoTrackDescriptor, bool)
1. SelectAudioTrackByAttribute.MarshalJSON() ([]byte, error)
1. SelectAudioTrackByID.AsAudioTrackDescriptor() (*AudioTrackDescriptor, bool)
1. SelectAudioTrackByID.AsBasicAudioTrackDescriptor() (BasicAudioTrackDescriptor, bool)
1. SelectAudioTrackByID.AsBasicTrackDescriptor() (BasicTrackDescriptor, bool)
1. SelectAudioTrackByID.AsBasicVideoTrackDescriptor() (BasicVideoTrackDescriptor, bool)
1. SelectAudioTrackByID.AsSelectAudioTrackByAttribute() (*SelectAudioTrackByAttribute, bool)
1. SelectAudioTrackByID.AsSelectAudioTrackByID() (*SelectAudioTrackByID, bool)
1. SelectAudioTrackByID.AsSelectVideoTrackByAttribute() (*SelectVideoTrackByAttribute, bool)
1. SelectAudioTrackByID.AsSelectVideoTrackByID() (*SelectVideoTrackByID, bool)
1. SelectAudioTrackByID.AsTrackDescriptor() (*TrackDescriptor, bool)
1. SelectAudioTrackByID.AsVideoTrackDescriptor() (*VideoTrackDescriptor, bool)
1. SelectAudioTrackByID.MarshalJSON() ([]byte, error)
1. SelectVideoTrackByAttribute.AsAudioTrackDescriptor() (*AudioTrackDescriptor, bool)
1. SelectVideoTrackByAttribute.AsBasicAudioTrackDescriptor() (BasicAudioTrackDescriptor, bool)
1. SelectVideoTrackByAttribute.AsBasicTrackDescriptor() (BasicTrackDescriptor, bool)
1. SelectVideoTrackByAttribute.AsBasicVideoTrackDescriptor() (BasicVideoTrackDescriptor, bool)
1. SelectVideoTrackByAttribute.AsSelectAudioTrackByAttribute() (*SelectAudioTrackByAttribute, bool)
1. SelectVideoTrackByAttribute.AsSelectAudioTrackByID() (*SelectAudioTrackByID, bool)
1. SelectVideoTrackByAttribute.AsSelectVideoTrackByAttribute() (*SelectVideoTrackByAttribute, bool)
1. SelectVideoTrackByAttribute.AsSelectVideoTrackByID() (*SelectVideoTrackByID, bool)
1. SelectVideoTrackByAttribute.AsTrackDescriptor() (*TrackDescriptor, bool)
1. SelectVideoTrackByAttribute.AsVideoTrackDescriptor() (*VideoTrackDescriptor, bool)
1. SelectVideoTrackByAttribute.MarshalJSON() ([]byte, error)
1. SelectVideoTrackByID.AsAudioTrackDescriptor() (*AudioTrackDescriptor, bool)
1. SelectVideoTrackByID.AsBasicAudioTrackDescriptor() (BasicAudioTrackDescriptor, bool)
1. SelectVideoTrackByID.AsBasicTrackDescriptor() (BasicTrackDescriptor, bool)
1. SelectVideoTrackByID.AsBasicVideoTrackDescriptor() (BasicVideoTrackDescriptor, bool)
1. SelectVideoTrackByID.AsSelectAudioTrackByAttribute() (*SelectAudioTrackByAttribute, bool)
1. SelectVideoTrackByID.AsSelectAudioTrackByID() (*SelectAudioTrackByID, bool)
1. SelectVideoTrackByID.AsSelectVideoTrackByAttribute() (*SelectVideoTrackByAttribute, bool)
1. SelectVideoTrackByID.AsSelectVideoTrackByID() (*SelectVideoTrackByID, bool)
1. SelectVideoTrackByID.AsTrackDescriptor() (*TrackDescriptor, bool)
1. SelectVideoTrackByID.AsVideoTrackDescriptor() (*VideoTrackDescriptor, bool)
1. SelectVideoTrackByID.MarshalJSON() ([]byte, error)
1. TrackDescriptor.AsAudioTrackDescriptor() (*AudioTrackDescriptor, bool)
1. TrackDescriptor.AsBasicAudioTrackDescriptor() (BasicAudioTrackDescriptor, bool)
1. TrackDescriptor.AsBasicTrackDescriptor() (BasicTrackDescriptor, bool)
1. TrackDescriptor.AsBasicVideoTrackDescriptor() (BasicVideoTrackDescriptor, bool)
1. TrackDescriptor.AsSelectAudioTrackByAttribute() (*SelectAudioTrackByAttribute, bool)
1. TrackDescriptor.AsSelectAudioTrackByID() (*SelectAudioTrackByID, bool)
1. TrackDescriptor.AsSelectVideoTrackByAttribute() (*SelectVideoTrackByAttribute, bool)
1. TrackDescriptor.AsSelectVideoTrackByID() (*SelectVideoTrackByID, bool)
1. TrackDescriptor.AsTrackDescriptor() (*TrackDescriptor, bool)
1. TrackDescriptor.AsVideoTrackDescriptor() (*VideoTrackDescriptor, bool)
1. TrackDescriptor.MarshalJSON() ([]byte, error)
1. Video.AsH265Video() (*H265Video, bool)
1. VideoLayer.AsBasicH265VideoLayer() (BasicH265VideoLayer, bool)
1. VideoLayer.AsH265Layer() (*H265Layer, bool)
1. VideoLayer.AsH265VideoLayer() (*H265VideoLayer, bool)
1. VideoTrackDescriptor.AsAudioTrackDescriptor() (*AudioTrackDescriptor, bool)
1. VideoTrackDescriptor.AsBasicAudioTrackDescriptor() (BasicAudioTrackDescriptor, bool)
1. VideoTrackDescriptor.AsBasicTrackDescriptor() (BasicTrackDescriptor, bool)
1. VideoTrackDescriptor.AsBasicVideoTrackDescriptor() (BasicVideoTrackDescriptor, bool)
1. VideoTrackDescriptor.AsSelectAudioTrackByAttribute() (*SelectAudioTrackByAttribute, bool)
1. VideoTrackDescriptor.AsSelectAudioTrackByID() (*SelectAudioTrackByID, bool)
1. VideoTrackDescriptor.AsSelectVideoTrackByAttribute() (*SelectVideoTrackByAttribute, bool)
1. VideoTrackDescriptor.AsSelectVideoTrackByID() (*SelectVideoTrackByID, bool)
1. VideoTrackDescriptor.AsTrackDescriptor() (*TrackDescriptor, bool)
1. VideoTrackDescriptor.AsVideoTrackDescriptor() (*VideoTrackDescriptor, bool)
1. VideoTrackDescriptor.MarshalJSON() ([]byte, error)

## Struct Changes

### New Structs

1. AudioTrackDescriptor
1. FromAllInputFile
1. FromEachInputFile
1. H265Layer
1. H265Video
1. H265VideoLayer
1. InputDefinition
1. InputFile
1. JobInputSequence
1. SelectAudioTrackByAttribute
1. SelectAudioTrackByID
1. SelectVideoTrackByAttribute
1. SelectVideoTrackByID
1. SystemData
1. TrackDescriptor
1. VideoTrackDescriptor

### New Struct Fields

1. AccountFilter.SystemData
1. Asset.SystemData
1. AssetFilter.SystemData
1. ContentKeyPolicy.SystemData
1. FaceDetectorPreset.BlurType
1. FaceDetectorPreset.Mode
1. Job.SystemData
1. JobInputAsset.InputDefinitions
1. JobInputClip.InputDefinitions
1. JobInputHTTP.InputDefinitions
1. LiveEvent.SystemData
1. MetricSpecification.LockAggregationType
1. Service.SystemData
1. StreamingEndpoint.SystemData
1. StreamingLocator.SystemData
1. StreamingPolicy.SystemData
1. Transform.SystemData
