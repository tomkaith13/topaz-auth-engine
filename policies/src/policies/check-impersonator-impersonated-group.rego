package policies.impersonator

default allowed = false

allowed {
	impersonator := ds.object({
		"object_id": input.user.id,
		"object_type": "user"
	})
	check := ds.check({
		"object_type": input.resource.object_type,
		"object_id": input.resource.object_id,
		"relation": input.resource.relation,
		"subject_type": "user",
		"subject_id": input.user.id,
	})
    impersonated_check := ds.check_relation({
        "object_type": "user",
        "object_id" : input.resource.impersonated_id,
        "relation": "impersonator",
        "subject_type": "user",
        "subject_id": impersonator.id

    })
    check
    impersonated_check
}