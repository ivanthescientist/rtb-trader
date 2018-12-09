# What is this? 

This is a proof of concept WIP minimalistic RTB DSP trader server. Right now it only contains 
few hardcoded non-overlapping segments (specifically for iOS, Android and Windows values of Device.OS field).
It responds with a $1e-15 bid to any BidRequest that belongs to one of those segments. 

Typically it uses fields from BidRequest, but it also contains few heuristics for extracting fields from UserAgent field also. 